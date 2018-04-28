package internals

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type RemoteGame struct {
	client *rpc.Client
}

func GetRemoteGameClient(serverAddress string) Game {
	client, err := rpc.DialHTTP("tcp", serverAddress) // WARN: address should include port
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return &RemoteGame{client}
}

func (rgame *RemoteGame) ApplyStep(step *StepArgs, ok *bool) error {
	return rgame.client.Call("RemoteGame.ApplyStep", step, ok)
}

func (rgame *RemoteGame) AddUser(user *Player, ok *bool) error {
	err := rgame.client.Call("RemoteGame.AddUser", user, ok)
	return err
}

func (rgame *RemoteGame) SetupGame(startArgs *GameStartArgs, reply *bool) error {
	// rgame.client.
	return rgame.client.Call("RemoteGame.SetupGame", startArgs, reply)
}

func (rgame *RemoteGame) Ping(value, reply *int) error {
	return rgame.client.Call("RemoteGame.Ping", value, reply)
}

func runRPCServer(game *RealGame, created chan<- bool) {

	server := rpc.NewServer()
	server.RegisterName("RemoteGame", game)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	port := GetRPCPort()

	game.rpcStopped = false
	game.rpcRunning = true

	for game.rpcRunning {
		l, e := net.Listen("tcp", port)
		game.rpcListener = &l
		//fmt.Println(l,e)
		if e != nil {
			log.Fatal("RPC: there was an error in listening for http connection on port "+port, e)
			created <- false
			return
		}
		log.Println("RPC: Started listening for new http connections.")
		created <- true
		err := http.Serve(l, nil)
		if err != nil {
			log.Println("RPC: Error serving connection.")
			continue
		}

		log.Println("RPC: Serving new connection.")
	}
	log.Printf("RPC: rpc server stopped")
	game.rpcStopped = true
}
