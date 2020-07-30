package types

import (
	"log"
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
	Conn        Conn
	Ticker      *time.Ticker
	ReadChan    chan []byte
	ReadBuffer  Buffer
	WriteBuffer Buffer
	Status      uint8
	OpPipe      *chan Op
}

func (s *Session) SetPipe(pipe *chan Op) {
	s.OpPipe = pipe
	s.Status = Pending
}

func (s *Session) Write(i []Op) (int, error) {
	data, err := s.WriteBuffer.Encode(i)
	if err != nil {
		log.Printf("Error encode data %v", err)
	}
	return s.Conn.Write(data)
}

func (s *Session) Read() {
	for {
		buffer, err := s.Conn.Read()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Printf("Error while Read msg %v", err)
				s.Status = Close
				data := make(map[string]float64)
				data["Id"] = float64(s.Id)
				op := Op{
					Type: Op_Logout,
					Data: data,
				}
				*s.OpPipe <- op
				return
			}
		} else {
			s.ReadChan <- buffer
		}
	}
}

func (s *Session) JoinRoom(roomId ID, ch *chan Op) {
	s.Status = InRoom
	s.OpPipe = ch
}
