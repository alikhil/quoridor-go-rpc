package internals

import (
	"github.com/googollee/go-socket.io"
	"log"
)

func CreateSocketServer(game *GGame) (*socketio.Server, error) {

	server, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}
	var res = new(bool) // this is useless, just for compatability with interface

	// server.
	server.On("connection", func(so socketio.Socket) {

		game.socket = &so
		// TODO:
		// We can reuse existing game to continue playing if connection was lost

		log.Println("SOCKET: Connection with frontend is established")

		so.On("create_game", func(userName string) {
			// TODO: think about repeating attempts to create game and connecting after creating own game etc
			log.Printf("SOCKET: create_game(%s) command recived", userName)
			game.StartSelfhostedGame()

			var selfHostedUser = Player{Name: &userName, Endpoint: GetIPAddress() + GetRPCPort()}
			game.AddUser(&selfHostedUser, res)
			so.Emit("show_endpoint", selfHostedUser.Endpoint)
		})

		so.On("connect_to_game", func(ip, name string) {
			if name == "" {
				name = "unnamed user"
			}
			log.Printf("SOCKET: connect_to_game(%s, %s) command recieved", ip, name)

			created := make(chan bool, 1)
			if !game.rpcRunning {
				go runRPCServer(game.RealGame, created)
				ok := <-created
				if !ok {
					log.Printf("SOCKET: failed to run rpc")
					return
				}
			}
			err := game.ConnectAsRemoteUser(&ConnectArgs{endpoint: ip, name: name}, res)
			if err != nil {
				log.Printf("SOCKET: failed to connect to game")
			}
			log.Printf("SOCKET: rpc call result was: %v", *res)
		})

		so.On("share_step", func(stepData StepData) {
			log.Printf("SOCKET: share step(%v) command recieved", stepData)
			if stepData.step != game.step {
				log.Printf("SOCKET: recived invalid step(%v); game step is %v", stepData.step, game.step)
			}
			game.step++
			err := game.ShareStep(stepData)
			if err != nil {
				log.Printf("SOCKET: failed to share step: %v", err)
			}
		})

		so.On("disconnection", func() {
			log.Println("SOCKET: connection with frontend is lost")

		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("SOCKET: error:", err)
	})
	return server, nil
}
