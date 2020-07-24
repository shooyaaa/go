package types

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	Close   = 1
	Open    = 2
	Pending = 4
	Waiting = 5
	InRoom  = 3
)

type Session struct {
	Id          ID
	Player      Player
	Conn        interface{}
	Ticker      *time.Ticker
	ReadChan    chan []byte
	ReadBuffer  Buffer
	WriteBuffer Buffer
	Status      uint8
	OpPipe      chan Op
}

func (s *Session) SetPipe(pipe chan Op) {
	s.OpPipe = pipe
	s.Status = Pending
}

func (s *Session) Write(i interface{}) error {
	data, _ := s.WriteBuffer.Encode(i)
	return s.Conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, data)
}

func (s *Session) JoinRoom(roomId ID, ch chan Op) {
	s.Status = InRoom
	s.OpPipe = ch
}
