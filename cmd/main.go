//usr/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	// "fmt"
	"log"
	"net/http"

	inter "github.com/alikhil/quoridor-go-rpc/internals"
)

const httpPort = 5000 // for socket, static files, http game start
const tcpPort = 5001  // for rpc over tcp

func main() {
	// log.Println(inter.GetIPAddress())
	// return
	game := inter.GGame{&inter.RealGame{}}

	server, err := inter.CreateSocketServer(&game)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../assets")))

	// http.HandleFunc("/connect-to-game")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil))

}

// server, err := CreateServer()
// if err != nil {
// 	log.Fatal(err)
// }

// http.Handle("/socket/", server)
// http.HandleFunc("/game/", GameInitiationHandler)

// http.Handle("/", http.FileServer(http.Dir("../assets")))

// log.Println(fmt.Sprintf("Serving at localhost:%v...", httpPort))
// log.Println("socket connection at /socket/")
// log.Println("static files at /")
// log.Println("http connection for game initiation at /game/")

// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil))
