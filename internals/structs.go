package internals

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	// "log"
)

const NeededPlayersCount = 2

type Player struct {
	Endpoint string
	Name     *string
}

type Game interface {
	AddUser(newUser *Player, ok *bool) error
	SetupGame(startArgs *GameStartArgs, ok *bool) error
}

type GameStartArgs struct {
	Players []Player
}

type RealGame struct {
	started       bool
	selfHosted    bool
	remoteAddress *string
	players       []Player
	step          int
	socket        *socketio.Socket
}

type GGame struct {
	*RealGame
}

type ConnectArgs struct {
	endpoint string
	name     string
}

func (game *GGame) StartSelfhostedGame() {
	game.selfHosted = true
	go runRPCServer(game.RealGame)
}

func EmitMakeStep(game *RealGame) {
	socket := *game.socket
	log.Printf("SOCKET: make_step emitted")
	socket.Emit("make_step")
}

func (player Player) IsHostedInThisMachine() bool {
	return player.Endpoint == GetEndpoint()
}

func (game *RealGame) AddUser(newUser *Player, ok *bool) error {
	for _, user := range game.players {
		if user == *newUser {
			return fmt.Errorf("Game: trying to add existing user(%v) to the game", newUser)
		}
	}
	game.players = append(game.players, *newUser)
	log.Printf("new user with name %s was added. now there is %v users", *newUser.Name, len(game.players))
	*ok = true
	if len(game.players) == NeededPlayersCount {
		startGame(game)
	}
	return nil
}

func (game *RealGame) ConnectAsRemoteUser(args *ConnectArgs, ok *bool) error {
	remoteGame := GetRemoteGameClient(args.endpoint)
	err := remoteGame.AddUser(&Player{Endpoint: GetIPAddress() + GetRPCPort(), Name: &args.name}, ok)
	if err != nil {
		log.Printf("failed to connect as user: %v", err)
	}
	return err
}

func (game *RealGame) SetupGame(startArgs *GameStartArgs, reply *bool) error {
	game.players = startArgs.Players
	*reply = true
	log.Printf("RPC: recived game start: %v", startArgs)
	go checkCurrentPlayer(game)
	return nil
}

func checkCurrentPlayer(game *RealGame) {
	for i, player := range game.players {
		if player.IsHostedInThisMachine() && game.step == i {
			EmitMakeStep(game)
			break
		}
	}
}

func startGame(game *RealGame) {
	// TODO: we can add shuffling here
	for _, player := range game.players {
		if player.IsHostedInThisMachine() {
			continue
		}
		ok := new(bool)
		remoteGame := GetRemoteGameClient(player.Endpoint)
		err := remoteGame.SetupGame(&GameStartArgs{Players: game.players}, ok)

		if err != nil {
			log.Printf("RPC: Failed to setup game for player %s with endpoint %s: %v", *player.Name, player.Endpoint, err)
		}
	}
	EmitMakeStep(game)
}
