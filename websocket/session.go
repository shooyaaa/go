package websocket

import (
	"github.com/gorilla/websocket"
)

type Session struct {
	Id   int64
	Name string
	Conn *websocket.Conn
}
