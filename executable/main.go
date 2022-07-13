package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"syscall"

	"github.com/shooyaaa/core/library"
	network2 "github.com/shooyaaa/core/network"
	types2 "github.com/shooyaaa/core/types"
	"github.com/shooyaaa/log"
	"github.com/shooyaaa/runnable/cron"
	"github.com/shooyaaa/runnable/env"
	"golang.org/x/net/ipv4"
)

func main() {

	i, _ := net.Interfaces()
	fmt.Println((i))
	fd, err := syscall.Socket(syscall.AF_ROUTE, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		fmt.Println("error while create raw socket ", err)
	}
	dmac, _ := net.ParseMAC("9c2b.a6dd.da41")
	smac, _ := net.ParseMAC("14:7d:da:da:3b:3a")
	arp, _ := network2.NewPacket(network2.OperationRequest, dmac, netip.MustParseAddr("127.0.0.1"), smac, netip.MustParseAddr("0.0.0.0"))
	b, _ := arp.MarshalBinary()
	network2.SendRaw(b)
	//err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		panic(err)
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		buf := make([]byte, 1500)
		f.Read(buf)
		ip4header, _ := ipv4.ParseHeader(buf)
		fmt.Println("ip header:", ip4header)
		switch ip4header.Protocol {
		case network2.IGMP_PROTOCOL:
			igmp := network2.NewIGMPHeader(buf[ip4header.Len:])
			fmt.Println("igmp ", igmp)
		case network2.TCP_PROTOCOL:
			tcpheader := network2.NewTCPHeader(buf[20:40])
			fmt.Println("tcp header: ", tcpheader)
		case network2.ICMP_PROTOCOL:
			icmp := network2.NewICMPHeader(buf[ip4header.Len:])
			fmt.Println(icmp)
		}
	}
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
