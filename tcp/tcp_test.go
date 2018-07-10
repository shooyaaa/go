package tcp

import (
	"github.com/shooyaaa/session"
	"github.com/shooyaaa/uuid"
	"testing"
)

func TestMain(t *testing.T) {
	simple := uuid.Simple{0}
	var uuid uuid.UUID
	uuid = &simple
	tcp := Tcp{
		Id:        uuid,
		Sessions:  make(map[int64]session.Session),
		HeartBeat: 5,
	}
	tcp.Listen("127.0.0.1:3333")
}