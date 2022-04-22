package main

import (
	"github.com/shooyaaa/core/codec"
	network2 "github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/session"
	types2 "github.com/shooyaaa/core/types"
	"github.com/shooyaaa/log"
	"github.com/shooyaaa/manager"
)

func main() {
	session.SetCodec(&codec.Json{})
	manager.SessionManager.Work()
	ws := network2.Ws{
		Id:        &types2.Simple{},
		HeartBeat: 40000000000,
		Root:      "./static",
	}
	go run(&ws, "127.0.0.1:5233")
	tcp := network2.Tcp{Id: &types2.Simple{}}
	run(&tcp, "127.0.0.1:3352")
}

func run(s network2.Server, addr string) {
	s.Listen(addr)
	log.Info("server listening on: %v", addr)
	for {
		manager.SessionManager.WaitChan <- s.Accept()
	}
	defer s.Close()
}
