package network

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"
)

/*
/*
#include <stdint.h>
#include <stdlib.h>

typedef struct __attribute__((packed))

	{
	    char dest[6];
	    char sender[6];
	    uint16_t protocolType;
	} EthernetHeader;

typedef struct __attribute__((packed))

	{
	    uint16_t hwType;
	    uint16_t protoType;
	    char hwLen;
	    char protocolLen;
	    uint16_t oper;
	    char SHA[6];
	    char SPA[4];
	    char THA[6];
	    char TPA[4];
	} ArpPacket;

typedef struct __attribute__((packed))

	{
	    EthernetHeader eth;
	    ArpPacket arp;
	} EthernetArpPacket;

char* FillRequestPacketFields(char* senderMac, char* senderIp)

	{
	    EthernetArpPacket * packet = malloc(sizeof(EthernetArpPacket));
	    memset(packet, 0, sizeof(EthernetArpPacket));
	    // Ethernet header
	    // Dest = Broadcast (ff:ff:ff:ff:ff)
	    packet->eth.dest[0] = 0xff;
	    packet->eth.dest[1] = 0xff;
	    packet->eth.dest[2] = 0xff;
	    packet->eth.dest[3] = 0xff;
	    packet->eth.dest[4] = 0xff;
	    packet->eth.dest[5] = 0xff;

	    packet->eth.sender[0] = strtol(senderMac, NULL, 16); senderMac += 3;
	    packet->eth.sender[1] = strtol(senderMac, NULL, 16); senderMac += 3;
	    packet->eth.sender[2] = strtol(senderMac, NULL, 16); senderMac += 3;
	    packet->eth.sender[3] = strtol(senderMac, NULL, 16); senderMac += 3;
	    packet->eth.sender[4] = strtol(senderMac, NULL, 16); senderMac += 3;
	    packet->eth.sender[5] = strtol(senderMac, NULL, 16);

	    packet->eth.protocolType = htons(0x0806); // ARP

	    // ARP Packet fields
	    packet->arp.hwType = htons(1); // Ethernet
	    packet->arp.protoType = htons(0x800); //IP;
	    packet->arp.hwLen = 6;
	    packet->arp.protocolLen = 4;
	    packet->arp.oper = htons(2); // response

	    // Sender MAC (same as that in eth header)
	    memcpy(packet->arp.SHA, packet->eth.sender, 6);

	    // Sender IP
	    packet->arp.SPA[0] = strtol(senderIp, NULL, 10); senderIp = strchr(senderIp, '.') + 1;
	    packet->arp.SPA[1] = strtol(senderIp, NULL, 10); senderIp = strchr(senderIp, '.') + 1;
	    packet->arp.SPA[2] = strtol(senderIp, NULL, 10); senderIp = strchr(senderIp, '.') + 1;
	    packet->arp.SPA[3] = strtol(senderIp, NULL, 10);

	    // Dest MAC: Same as SHA, as we use an ARP response
	    memcpy(packet->arp.THA, packet->arp.SHA, 6);

	    // Dest IP: Same as SPA
	    memcpy(packet->arp.TPA, packet->arp.SPA, 4);

	    return (char*) packet;
	}

import "C"

	func SendArp1() {
		etherArp := new(C.EthernetArpPacket)
		size := uint(unsafe.Sizeof(*etherArp))
		fmt.Println("Size : ", size)

		fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		fmt.Println("Obtained fd ", fd)
		defer syscall.Close(fd)

		// Get Mac address of vboxnet1
		interf, err := net.InterfaceByName("enp4s0")
		if err != nil {
			fmt.Println("Could not find vboxnet interface")
			return
		}

		fmt.Println("Interface hw address: ", interf.HardwareAddr)
		fmt.Println("Creating request for IP 10.10.10.2 from IP 10.10.10.1")

		iface_cstr := C.CString(interf.HardwareAddr.String())
		ip_cstr := C.CString("10.10.10.5")

		packet := C.GoBytes(unsafe.Pointer(C.FillRequestPacketFields(iface_cstr, ip_cstr)), C.int(size))

		// Send the packet
		var addr syscall.SockaddrLinklayer
		addr.Protocol = syscall.ETH_P_ARP
		addr.Ifindex = interf.Index
		addr.Hatype = syscall.ARPHRD_ETHER

		//err = syscall.Sendto(fd, packet, 0, &addr)
		ifc, _ := MainInterface()
		err = SendRaw(packet, ifc.Name)

		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Sent packet")
		}

}
*/
const (
	ICMP_PROTOCOL = 1
	IGMP_PROTOCOL = 2
	TCP_PROTOCOL  = 6
	UDP_PROTOCOL  = 17

	ICMP_REQUEST     = 8
	ICMP_REPLY       = 0
	ICMP_UNREACHABLE = 3
)

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}

type IGMP struct {
	Type         uint8
	MaxRespTime  uint8
	Checksum     uint16
	GroupAddress int
}

func (igmp IGMP) String() string {
	msgType := "Unknow"
	if igmp.Type == 0x11 {
		msgType = "Membership Query"
	} else if igmp.Type == 0x12 {
		msgType = "IGMPv1 Memmbership Report"
	} else if igmp.Type == 0x16 {
		msgType = "IGMPv2 Membership Report"
	} else if igmp.Type == 0x22 {
		msgType = "IGMPv2 Membership Report"
	} else if igmp.Type == 0x17 {
		msgType = "Leave Group"
	}
	return fmt.Sprintf("IGMP msg type %s, Max Resp Time %v, Group Address %v", msgType, igmp.MaxRespTime, igmp.GroupAddress)
}

func NewIGMPHeader(bts []byte) IGMP {
	igmp := IGMP{}
	buf := bytes.NewBuffer(bts)
	binary.Read(buf, binary.BigEndian, &igmp.Type)
	binary.Read(buf, binary.BigEndian, &igmp.MaxRespTime)
	binary.Read(buf, binary.BigEndian, &igmp.GroupAddress)
	return igmp
}

type ICMP struct {
	Type         uint8
	Code         uint8
	Checksum     uint16
	RestOfHeader uint32
	Data         []byte
}

func (i *ICMP) CalCheckSum() {
	data := i.Marshal()
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	i.Checksum = uint16(^sum)
}

func (i ICMP) IsUnreachable() bool {
	return i.Type == ICMP_UNREACHABLE
}

func (i ICMP) UnreachableHost() net.IP {
	if i.Type == ICMP_UNREACHABLE {
		ipHeader, err := ipv4.ParseHeader(i.Data)
		if err == nil {
			return ipHeader.Dst
		}
	}
	return nil
}

func (i ICMP) Marshal() []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, i.Type)
	binary.Write(&buffer, binary.BigEndian, i.Code)
	binary.Write(&buffer, binary.BigEndian, i.Checksum)
	binary.Write(&buffer, binary.BigEndian, i.RestOfHeader)
	if len(i.Data) > 0 {
		binary.Write(&buffer, binary.BigEndian, i.Data)
	}
	return buffer.Bytes()
}

func (icmp ICMP) String() string {
	msgType := "Unknow"
	if icmp.Type == ICMP_REPLY {
		msgType = "Echo Reply"
	} else if icmp.Type == 3 {
		msgType = []string{
			"Destination network unreachable",
			"Destination host unreachable",
			"Destination protocol unreachable",
			"Destination port unreachable",
			"Fragmentation required, and DF flag set",
			"Source route failed",
			"Destination network unknown",
			"Destination host unknown",
			"Source host isolated",
			"Network administratively prohibited",
			"Host administratively prohibited",
			"Network unreachable for ToS",
			"Host unreachable for ToS",
			"Communication administratively prohibited",
			"Host Precedence Violation",
			"Precedence cutoff in effect",
		}[icmp.Code]
	} else if icmp.Type == ICMP_REQUEST {
		msgType = "Echo Request"
	} else if icmp.Type == 11 {
		if icmp.Code == 0 {
			msgType = "TTL expired in transit"
		} else {
			msgType = "Fragment reassembly time exceeded"
		}
	} else if icmp.Type == 30 {
		msgType = "Traceroute"
	}
	return fmt.Sprintf("ICMP type %s", msgType)
}

func NewICMP(bts []byte) ICMP {
	icmp := ICMP{}
	buf := bytes.NewBuffer(bts)
	binary.Read(buf, binary.BigEndian, &icmp.Type)
	binary.Read(buf, binary.BigEndian, &icmp.Code)
	binary.Read(buf, binary.BigEndian, &icmp.Checksum)
	binary.Read(buf, binary.BigEndian, &icmp.RestOfHeader)
	icmp.Data = make([]byte, buf.Len())
	binary.Read(buf, binary.BigEndian, &icmp.Data)
	return icmp
}

type TCPHeader struct {
	Source      uint16
	Destination uint16
	SeqNum      uint32
	AckNum      uint32
	DataOffset  uint8 // 4 bits
	Reserved    uint8 // 3 bits
	ECN         uint8 // 3 bits
	Ctrl        uint8 // 6 bits
	Window      uint16
	Checksum    uint16 // Kernel will set this if it's 0
	Urgent      uint16
	Options     []TCPOption
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

func NewTCPHeader(bts []byte) TCPHeader {
	th := TCPHeader{}
	buf := bytes.NewBuffer(bts)
	binary.Read(buf, binary.BigEndian, &th.Source)
	binary.Read(buf, binary.BigEndian, &th.Destination)
	binary.Read(buf, binary.BigEndian, &th.SeqNum)
	binary.Read(buf, binary.BigEndian, &th.AckNum)
	binary.Read(buf, binary.BigEndian, &th.DataOffset)
	binary.Read(buf, binary.BigEndian, &th.Reserved)
	binary.Read(buf, binary.BigEndian, &th.ECN)
	binary.Read(buf, binary.BigEndian, &th.Ctrl)
	binary.Read(buf, binary.BigEndian, &th.Window)
	binary.Read(buf, binary.BigEndian, &th.Checksum)
	binary.Read(buf, binary.BigEndian, &th.Urgent)
	binary.Read(buf, binary.BigEndian, &th.Options)
	return th
}

func (tcp *TCPHeader) Marshal() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, tcp.Source)
	binary.Write(buf, binary.BigEndian, tcp.Destination)
	binary.Write(buf, binary.BigEndian, tcp.SeqNum)
	binary.Write(buf, binary.BigEndian, tcp.AckNum)

	mix := uint16(tcp.DataOffset)<<12 | // top 4 bits
		uint16(tcp.Reserved)<<9 | // 3 bits
		uint16(tcp.ECN)<<6 | // 3 bits
		uint16(tcp.Ctrl) // bottom 6 bits
	binary.Write(buf, binary.BigEndian, mix)

	binary.Write(buf, binary.BigEndian, tcp.Window)
	binary.Write(buf, binary.BigEndian, tcp.Checksum)
	binary.Write(buf, binary.BigEndian, tcp.Urgent)

	for _, option := range tcp.Options {
		binary.Write(buf, binary.BigEndian, option.Kind)
		if option.Length > 1 {
			binary.Write(buf, binary.BigEndian, option.Length)
			binary.Write(buf, binary.BigEndian, option.Data)
		}
	}

	out := buf.Bytes()

	// Pad to min tcp header size, which is 20 bytes (5 32-bit words)
	pad := 20 - len(out)
	for i := 0; i < pad; i++ {
		out = append(out, 0)
	}

	return out
}

// TCP Checksum
func Csum(data []byte, srcip, dstip [4]byte) uint16 {

	pseudoHeader := []byte{
		srcip[0], srcip[1], srcip[2], srcip[3],
		dstip[0], dstip[1], dstip[2], dstip[3],
		0,                  // zero
		6,                  // router number (6 == TCP)
		0, byte(len(data)), // TCP length (16 bits), not inc pseudo header
	}

	sumThis := make([]byte, 0, len(pseudoHeader)+len(data))
	sumThis = append(sumThis, pseudoHeader...)
	sumThis = append(sumThis, data...)
	//fmt.Printf("% x\n", sumThis)

	lenSumThis := len(sumThis)
	var nextWord uint16
	var sum uint32
	for i := 0; i+1 < lenSumThis; i += 2 {
		nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
		sum += uint32(nextWord)
	}
	if lenSumThis%2 != 0 {
		//fmt.Println("Odd byte")
		sum += uint32(sumThis[len(sumThis)-1])
	}

	// Add back any carry, and any carry from adding the carry
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// Bitwise complement
	return uint16(^sum)
}

var (
	ifName  = flag.String("i", "", "Interface Name")
	portNum = flag.Uint("p", 0, "Port to inspect")
	host    = flag.String("h", "", "Host to send")
)

func ifAddr(ifName string) net.Addr {
	info, err := net.InterfaceByName(ifName)
	if err != nil {
		log.Fatalf("net.InterfaceByName for %s. %s", ifName, err)
	}

	addrs, err := info.Addrs()
	if err != nil {
		log.Fatalf("iface.Addrs: %s", err)
	}
	return addrs[0]
}

func sendSyn(sAddr, dAddr string, port uint16) error {
	packet := TCPHeader{
		Source:      0xaa47, // Random ephemeral port
		Destination: port,
		SeqNum:      rand.Uint32(),
		AckNum:      0,
		DataOffset:  5,      // 4 bits
		Reserved:    0,      // 3 bits
		ECN:         0,      // 3 bits
		Ctrl:        2,      // 6 bits (000010, SYN bit set)
		Window:      0xaaaa, // The amount of data that it is able to accept in bytes
		Checksum:    0,      // Kernel will set this if it's 0
		Urgent:      0,
		Options:     []TCPOption{},
	}

	data := packet.Marshal()
	packet.Checksum = Csum(data, to4byte(sAddr), to4byte(dAddr))

	data = packet.Marshal()

	fmt.Printf("% x\n", data)

	conn, err := net.Dial("ip4:tcp", dAddr)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	numWrote, err := conn.Write(data)
	if err != nil {
		log.Fatalf("Write: %s\n", err)
	}
	if numWrote != len(data) {
		log.Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
	}

	conn.Close()

	return nil
}

func SendRaw(p []byte, ifcName string) error {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		log.Fatal("failed to create raw scoket ", err)
	}
	interf, _ := net.InterfaceByName(ifcName)
	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ARP
	addr.Ifindex = interf.Index
	addr.Hatype = syscall.ARPHRD_ETHER
	err = syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	}
	return err
}

func Ping(target [4]byte, timeout time.Duration) bool {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Println("failed to create raw socket ", err)
		return false
	}
	icmp := ICMP{Type: ICMP_REQUEST}

	var addr syscall.SockaddrInet4
	addr.Port = 0
	addr.Addr = target
	dst := net.IPv4(target[0], target[1], target[2], target[3])
	ip := ipv4.Header{Version: 4, Len: 20, TotalLen: 30, TTL: 64, Protocol: 1, Dst: dst}
	bytes, _ := ip.Marshal()
	icmp.CalCheckSum()
	err = syscall.Sendto(fd, append(bytes, icmp.Marshal()...), 0, &addr)
	if err != nil {
		log.Println("failed to send icmp package to host " + ip.String() + err.Error())
		return false
	}
	ch := make(chan bool, 1)
	go CaptureIcmp(dst, timeout, ch)
	for {
		select {
		case ret := <-ch:
			return ret
		case <-time.After(timeout):
			return false
		}
	}
}

func CaptureIcmp(src net.IP, duration time.Duration, ch chan bool) {
	fd1, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	f := os.NewFile(uintptr(fd1), fmt.Sprintf("fd %d", fd1))
	defer f.Close()
	buf := make([]byte, 1500)
	for {
		f.Read(buf)
		ip4header, _ := ipv4.ParseHeader(buf)
		switch ip4header.Protocol {
		case ICMP_PROTOCOL:
			icmp := NewICMP(buf[ip4header.Len:])
			if icmp.IsUnreachable() && icmp.UnreachableHost().Equal(src) {
				ch <- false
				return
			}
			if ip4header.Src.Equal(src) && icmp.Type == ICMP_REPLY {
				ch <- true
				return
			}
		}
	}
}
