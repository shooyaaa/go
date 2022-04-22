package network

import (
	"github.com/shooyaaa/core/session"
	types2 "github.com/shooyaaa/core/types"
	"testing"
)

func TestMain(m *testing.M) {
	simple := types2.Simple{}
	var uuid types2.UUID
	uuid = &simple
	tcp := Tcp{
		Id:        uuid,
		Sessions:  make(map[int64]session.Session),
		HeartBeat: 5,
	}
	tcp.Listen("127.0.0.1:3333")
}
