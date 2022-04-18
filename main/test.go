package main

import (
	"github.com/shooyaaa/log"
	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/network"
	"github.com/shooyaaa/types"
)

func main() {
	types.SetCodec(&types.Json{})
	manager.SessionManager().Work()
	ws := network.Ws{
		Id:        &types.Simple{},
		HeartBeat: 40000000000,
		Root:      "./static",
	}
	go run(&ws, "127.0.0.1:5233")
	tcp := network.Tcp{Id: &types.Simple{}}
	run(&tcp, "127.0.0.1:3352")
}

func run(s network.Server, addr string) {
	s.Listen(addr)
	log.Info("server listening on: %v", addr)
	for {
		manager.SessionManager().WaitChan <- s.Accept()
	}
	defer s.Close()
}
