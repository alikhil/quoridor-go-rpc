package internals

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net"
	"time"
	// "log"
)

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
	StepID  int
}

type RealGame struct {
	started                bool
	selfHosted             bool
	players                []Player
	step                   int
	socket                 *socketio.Socket
	rpcRunning             bool
	rpcStopped             bool
	clients                map[string]Game
	ticker                 *time.Ticker
	numberOfPlayers        int
	healthcheckerIsRunning bool
	rpcListener            *net.Listener
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
	Data map[string]int
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
	return err == nil && *a == *b
}

func stop(game *RealGame) {
	game.ticker.Stop()
	game.rpcRunning = false
	(*game.rpcListener).Close()
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
	if len(game.players) == game.numberOfPlayers {
		*ok = false
		return nil
	}

	for _, user := range game.players {
		if user == *newUser {
			return fmt.Errorf("Game: trying to add existing user(%v) to the game", newUser)
		}
	}
	game.players = append(game.players, *newUser)
	log.Printf("new user with name %s was added. now there is %v users", *newUser.Name, len(game.players))
	*ok = true
	if len(game.players) == game.numberOfPlayers {
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
	game.step = startArgs.StepID
	*reply = true
	log.Printf("RPC: recived game start: %v", startArgs)
	go checkCurrentPlayer(game)
	go startHealchecker(game)
	if !game.started {
		(*game.socket).Emit("on_create", len(startArgs.Players))
		game.started = true
	}
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

	if game.healthcheckerIsRunning {
		log.Printf("HELATH: hey healthchecker is already runnung")
		return
	}
	log.Printf("HEALTH: hey! healthchecker works!")
	game.ticker = time.NewTicker(time.Millisecond * 1000)
	game.healthcheckerIsRunning = true
	// game.ticker.
	for _ = range game.ticker.C {

		// log.Printf("Start healthcheck %v", t)
		var playersCnt = len(game.players)
		var statuses = make([]bool, playersCnt)
		var thereIsFailed = false
		var alivePlayersCnt = 1
		var newPlayers []Player

		curPlayerI := 0

		for i, player := range game.players {
			if player.IsHostedInThisMachine() {
				curPlayerI = i
				statuses[i] = true
				newPlayers = append(newPlayers, player)
				continue
			}
			statuses[i] = isConnected(getRemoteGame(player.Endpoint, game))
			if statuses[i] {
				alivePlayersCnt++
				newPlayers = append(newPlayers, player)
			}
			thereIsFailed = thereIsFailed || !statuses[i]
		}

		if thereIsFailed {
			log.Printf("HEALTH: there is failed node; status: %v", statuses)

			if alivePlayersCnt == 1 {
				log.Printf("HEALTH: wow I am the only lonely node that alive :( ; Nobody wants to play with me. Stopping the game server")
				stop(game)
				continue
			}

			steper := game.step % playersCnt
			firstActiveAfterFailed := -1
			failed := false
			// search for fixer node
			for i, status := range statuses {
				if i < steper { // we want to find first alive node after failed
					continue
				}
				if failed && status {
					firstActiveAfterFailed = i
					break
				}
				if !status {
					failed = true
				}
			}

			if firstActiveAfterFailed == -1 {
				// we did not found alive node at the end of list, there can be in the beggining
				for i, status := range statuses {
					if status {
						firstActiveAfterFailed = i
						break
					}
				}
			}

			if firstActiveAfterFailed == -1 {
				log.Printf("HEALTH: no recovery nodes found [some shit is happening here; call the game administrator]")
				continue
			}

			if firstActiveAfterFailed == curPlayerI {

				newStepID := -1
				for i, player := range newPlayers {
					if player.IsHostedInThisMachine() {
						newStepID = i
						break
					}
				}

				if newStepID == -1 {
					log.Printf("HEALTH: something unbelievable happened. there is no current player in list of new players")
					continue
				}

				setupSucceeded := false
				for _, player := range newPlayers {
					remoteGame := getRemoteGame(player.Endpoint, game)
					ok := new(bool)
					err := remoteGame.SetupGame(&GameStartArgs{Players: newPlayers, StepID: newStepID}, ok)
					if err != nil {
						log.Printf("HEALTH: failed to setup game for player (%v)", player)
						if setupSucceeded {
							log.Printf("HEALTH: failed to recover; there can be inconsistency problems; ")
							panic("RECOVERY FAILED")
						}
					} else {
						setupSucceeded = true
					}
				}
				game.step = newStepID
				game.players = newPlayers
				EmitMakeStep(game, newPlayers[newStepID].PawnID)
				// TODO: WARN:
				// we will re set game set
				// but it can have problems because of concurrency
				// if player made step before we fix everything

			} else {
				log.Printf("HEALTH: it is not my turn, so just ignoring problem")
			}
		}
		game.healthcheckerIsRunning = false
		// time.Sleep(time.NewTicker().
	}
}
