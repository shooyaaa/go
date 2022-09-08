package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/shooyaaa/core/library"
	network2 "github.com/shooyaaa/core/network"
	types2 "github.com/shooyaaa/core/types"
	"github.com/shooyaaa/log"
	"github.com/shooyaaa/runnable/cron"
	"github.com/shooyaaa/runnable/env"
)

func main() {
	ifc, _ := network2.MainInterface()
	hosts := network2.ScanHosts(ifc.Name)
	fmt.Println("hosts ", hosts)
	for _, h := range hosts {
		fmt.Println(net.LookupAddr(h.String()))
	}

	ws := network2.Ws{
		Id:        &types2.Simple{},
		HeartBeat: 40000000000,
		Root:      "./static",
	}
	go run(&ws, "127.0.0.1:5233")
	tcp := network2.Tcp{Id: &types2.Simple{}}
	go run(&tcp, "127.0.0.1:3352")
	cron := cron.Cron{}

	library.ModuleManager.Load(env.Env{}.Run())
	library.ModuleManager.Load(cron.Run())
	go library.ModuleManager.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	s := <-ch
	fmt.Println("signal caught", s)
	library.ModuleManager.Exit()
}

func run(s network2.Server, addr string) {
	s.Listen(addr)
	log.Info("server listening on: %v", addr)
	defer s.Close()
}
