package types

import (
	"time"
)

const (
	Close  = 1
	Open   = 2
	InRoom = 3
)

type Session struct {
	Id       	ID
	Player		Player
	Conn 		interface{}
	Ticker   	*time.Ticker
	ReadChan 	chan []byte
	ReadBuffer  Buffer
	WriteBuffer Buffer
	Status		uint8
	OpPipe		chan Op
}
