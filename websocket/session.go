package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Session struct {
	Id       int64
	Name     string
	Conn     *websocket.Conn
	Ticker   *time.Ticker
	ReadChan chan []byte
}

func (s *Session) closeHandler(code int, text string) error {
	s.Ticker.Stop()
	close(s.ReadChan)
	s.Conn.Close()
	return nil
}

func (s Session) Read(ch chan []byte) {
	for {
		if s.Conn == nil {
			break
		}
		_, message, err := s.Conn.ReadMessage()
		log.Println("read message")
		if err != nil {
			log.Println("Read websocket error :", err)
			break
		} else {
			ch <- message
		}
	}
}
