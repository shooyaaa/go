package network

import (
	"testing"

	"github.com/shooyaaa/types"
)

func TestMain(t *testing.T) {
	simple := types.Simple{}
	var uuid types.UUID
	uuid = &simple
	tcp := Tcp{
		Id:        uuid,
		Sessions:  make(map[types.ID]types.Session),
		HeartBeat: 5,
	}
	tcp.Listen("127.0.0.1:3333")
}
