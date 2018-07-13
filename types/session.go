package types

import (
	"github.com/shooyaaa/uuid"
	"time"
)

type Session struct {
	Id       uuid.ID
	Name     string
	Conn     interface{}
	Ticker   *time.Ticker
	ReadChan chan []byte
	Buffer   Buffer
}
