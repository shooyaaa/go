package main

import (
	"github.com/shooyaaa/http"
	"github.com/shooyaaa/uuid"
	"github.com/shooyaaa/websocket"
)

func main() {
	simple := uuid.Simple{0}
	var uuid uuid.UUID
	uuid = &simple
	ws := websocket.Ws{Id: uuid, Sessions: make(map[int64]websocket.Session), HeartBeat: 5}
	server := http.HttpServer{"./static/", "127.0.0.1:3333", make(map[string]http.HttpHandler)}
	server.Register("/ws", ws.Connect)
	server.Register("/wsinfo", server.Info)
	server.Run()
}
