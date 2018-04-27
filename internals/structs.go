package internals

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"time"
	// "log"
)

const NeededPlayersCount = 2

type Player struct {
	Endpoint string
	Name     *string
	PawnID   int
}

type Game interface {
	AddUser(newUser *Player, ok *bool) error
	SetupGame(startArgs *GameStartArgs, ok *bool) error
	ApplyStep(step *StepArgs, ok *bool) error
	Ping(val, reply *int) error
}

type GameStartArgs struct {
	Players []Player
}

type RealGame struct {
	started    bool
	selfHosted bool
	players    []Player
	step       int
	socket     *socketio.Socket
	rpcRunning bool
	rpcStopped bool
	clients    map[string]Game
	ticker     *time.Ticker
}

type GGame struct {
	*RealGame
}

type ConnectArgs struct {
	endpoint string
	name     string
}

type StepData struct {
	Step int
	Data string
}

type StepArgs struct {
	Data StepData
}

func CreateGame() *GGame {
	game := &GGame{&RealGame{}}
	game.clients = make(map[string]Game)
	return game
}

func getRemoteGame(endpoint string, game *RealGame) Game {
	if game.clients[endpoint] == nil {
		game.clients[endpoint] = GetRemoteGameClient(endpoint)
	}
	return game.clients[endpoint]
}

func isConnected(game Game) bool {
	a := new(int)
	b := new(int)
	*a = 12
	err := game.Ping(a, b)
	return err != nil && *a == *b
}

func (ggame *GGame) Stop() {
	ggame.ticker.Stop()
	ggame.rpcRunning = false
}

func (game *GGame) StartSelfhostedGame() {
	game.players = []Player{}
	game.selfHosted = true
	if game.rpcRunning {
		log.Printf("RPC is already runnig do nothing...")
		return
	}

	created := make(chan bool, 1)
	go runRPCServer(game.RealGame, created)
	log.Printf("RPC successfully started: %v", <-created)
}

func (game *GGame) ShareStep(step StepData) error {
	for _, player := range game.players {
		if player.IsHostedInThisMachine() {
			continue
		}
		remoteGame := getRemoteGame(player.Endpoint, game.RealGame)
		res := new(bool)
		err := remoteGame.ApplyStep(&StepArgs{Data: step}, res)
		if err != nil {
			return err
		}
	}
	return nil
}

func EmitMakeStep(game *RealGame, index int) {
	log.Printf("SOCKET: make_step emitted")
	(*game.socket).Emit("make_step", game.step, index)
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

func (game *RealGame) Ping(value, reply *int) error {
	*reply = *value
	return nil
}

func (game *RealGame) ConnectAsRemoteUser(args *ConnectArgs, ok *bool) error {
	remoteGame := getRemoteGame(args.endpoint, game)
	err := remoteGame.AddUser(&Player{Endpoint: GetIPAddress() + GetRPCPort(), Name: &args.name}, ok)
	if err != nil {
		log.Printf("failed to connect as user: %v", err)
	}
	return err
}

func (game *RealGame) ApplyStep(step *StepArgs, ok *bool) error {

	game.step++
	(*game.socket).Emit("apply_step", step.Data)
	log.Printf("SOCKET: apply_step(%v)", step)
	go checkCurrentPlayer(game)
	*ok = true
	return nil
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
		if player.IsHostedInThisMachine() && game.step%len(game.players) == i {
			EmitMakeStep(game, player.PawnID)
			break
		}
	}
}

func startGame(game *RealGame) {
	// TODO: we can add shuffling here
	for i, player := range game.players {
		player.PawnID = i
		if player.IsHostedInThisMachine() {
			continue
		}
		ok := new(bool)
		remoteGame := getRemoteGame(player.Endpoint, game)
		err := remoteGame.SetupGame(&GameStartArgs{Players: game.players}, ok)

		if err != nil {
			// TODO: and what to do with such player?
			log.Printf("RPC: Failed to setup game for player %s with endpoint %s: %v", *player.Name, player.Endpoint, err)
		}
	}
	// assuming that we are first in list of players
	EmitMakeStep(game, 0)
	go startHealchecker(game)
}

func startHealchecker(game *RealGame) {

	game.ticker = time.NewTicker(time.Millisecond * 200)
	
	
	for _ = range game.ticker.C {
		
		var playersCnt = len(game.players)
		var statuses = make([]bool, playersCnt)
		var thereIsFailed = false
		for i, player := range game.players {
			if player.IsHostedInThisMachine() {
				continue
			}
			statuses[i] = isConnected(getRemoteGame(player.Endpoint, game))
			thereIsFailed = thereIsFailed || statuses[i]
		}

		if thereIsFailed {
			log.Printf("HEALTH: there is failed node; status: %v", statuses)

			if game.step % playersCnt == 
		}
		// time.Sleep(time.NewTicker().
	}
}
