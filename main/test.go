package main

import (
	"github.com/shooyaaa/connector"
	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/types"
)

func main() {
	ws := connector.Ws{
		Id:        &types.Simple{},
		HeartBeat: 40000000000,
		Addr:      "127.0.0.1:5233",
		Root:      "../static",
	}
	manager.SessionManager().Work()
	go ws.Run()
	simple := types.Simple{}
	var uuid types.UUID
	uuid = &simple
	tcp := connector.Tcp{
		Id:        uuid,
		Sessions:  make(map[types.ID]types.Session),
		HeartBeat: 5,
	}
	tcp.Listen("127.0.0.1:3352")
}
