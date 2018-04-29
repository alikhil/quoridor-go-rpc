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

func main() {
	game := inter.CreateGame()

	server, err := inter.CreateSocketServer(game)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../assets")))

	// http.HandleFunc("/connect-to-game")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil))

}
