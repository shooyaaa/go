package main

import (
	"fmt"
	"net"
	"net/netip"
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
	hosts := network2.ScanHosts("enp4s0")
	fmt.Println("hosts ", hosts)

	dmac, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff")
	smac, _ := net.ParseMAC("14:7d:da:da:3b:3a")
	arp, _ := network2.NewPacket(network2.OperationRequest, dmac, netip.MustParseAddr("127.0.0.1"), smac, netip.MustParseAddr("0.0.0.0"))
	b, _ := arp.MarshalBinary()
	network2.SendRaw(b)
	//err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	/*addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}
	p := pkt()
	err = syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	}*/
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
