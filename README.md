# quoridor-go-rpc

Distributed Quoridor game - implemented using Go and Javascript


## Quickstart

```bash
cd ~/go/src/github
git clone git@github.com:alikhil/quoridor-go-rpc.git alikhil/quoridor-go-rpc

go get ./...

cd cmd
go run main.go
```


## Components

* Frontend browser app with socket connection to backend

* Backend or p2p app - communicates with other peers and sends game state to Frontend

## Resources

* Sockets - [socket.io for go](https://github.com/googollee/go-socket.io)

* Quoridor clients - [example 1](https://github.com/danielborowski/quoridor-ai), [example 2](https://github.com/ranjez/Quoridor)

* RPC - [tutorial](https://medium.com/@akashg/remote-procedure-calls-with-go-1b85eb93b491), [docs](https://golang.org/pkg/net/rpc/), the book sent by proffessor.