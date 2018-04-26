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

func (rgame *RemoteGame) AddUser(user *Player, ok *bool) error {
	err := rgame.client.Call("RemoteGame.AddUser", user, ok)
	return err
}

func (rgame *RemoteGame) SetupGame(startArgs *GameStartArgs, reply *bool) error {
	return rgame.client.Call("RemoteGame.SetupGame", startArgs, reply)
}

func runRPCServer(game *RealGame) {

	server := rpc.NewServer()
	server.RegisterName("RemoteGame", game)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	port := GetRPCPort()

	for true {
		l, e := net.Listen("tcp", port)
		//fmt.Println(l,e)
		if e != nil {
			log.Fatal("RPC: there was an error in listening for http connection on port "+port, e)
			return
		} else {
			log.Println("RPC: Started listening for new http connections.")
		}

		err := http.Serve(l, nil)
		if err != nil {
			log.Println("RPC: Error serving connection.")
			continue
		}

		log.Println("RPC: Serving new connection.")
	}
}
