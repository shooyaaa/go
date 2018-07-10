package connector

import (
	"github.com/shooyaaa/codec"
	"github.com/shooyaaa/http"
	"github.com/shooyaaa/uuid"
	"github.com/shooyaaa/websocket"
)

func Run(root, addr) http.HttpServer {
	simple := uuid.Simple{0}
	var uuid uuid.UUID
	uuid = &simple
	ws := websocket.Ws{
		Id:        uuid,
		Sessions:  make(map[int64]websocket.Session),
		HeartBeat: 5,
		Codec:     &codec.Json{},
	}
	server := http.HttpServer{root, addr, make(map[string]http.HttpHandler)}
	server.Register("/ws", ws.Connect)
	server.Register("/wsinfo", server.Info)
	server.Run()
	return server
}
