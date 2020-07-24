package main

import (
	"github.com/shooyaaa/connector"
	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/types"
)

func main() {
	ws := connector.Ws{
		Id:        &types.Simple{},
		HeartBeat: 400,
		Addr:      "127.0.0.1:5233",
		Root:      "../static",
	}
	manager.SessionManager().Work()
	ws.Run()
}
