package network

import (
	"testing"

	types2 "github.com/shooyaaa/core/types"
)

func TestMain(m *testing.M) {
	simple := types2.Simple{}
	var uuid types2.UUID
	uuid = &simple
	tcp := Tcp{
		Id:        uuid,
		HeartBeat: 5,
	}
	tcp.Listen("127.0.0.1:3333")
}
