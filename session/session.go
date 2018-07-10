package session

import (
	"github.com/shooyaaa/codec"
	"time"
)

type Session struct {
	Id       int64
	Name     string
	Conn     interface{}
	Ticker   *time.Ticker
	ReadChan chan []byte
	Buffer   codec.Buffer
}

type Connection interface {
	Close()
}

func (s *Session) CloseHandler(code int, text string) error {
	s.Ticker.Stop()
	close(s.ReadChan)
	s.Conn.(Connection).Close()
	return nil
}
