package main

import (
	"github.com/shooyaaa/connector"
	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/types"
)

func main() {
	sm :=  manager.Session {
		WaitChan : make (chan types.Session, 1000),
	}
	sm.Init()
	ws := connector.Ws{
		Id : &types.Simple{},
		SessionManager : sm,
		HeartBeat : 400,
		Addr : "127.0.0.1:5233",
		Root : "../static",
	}
	sm.Work()
	ws.Run()
}
