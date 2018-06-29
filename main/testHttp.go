package main

import (
	"github.com/shooyaaa/connector"
	"github.com/shooyaaa/websocket"
)

func main() {
	ws := websocket.Ws{Addr: "127.0.0.1:8888"}
	server := connector.HttpServer{"./static/", "127.0.0.1:3333", make(map[string]connector.HttpHandler)}
	server.Register("/ws", ws.Connect)
	server.Register("/wsinfo", server.Info)
	server.Run()
}
