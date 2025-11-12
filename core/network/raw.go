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
	return nil
}

func Ping(target [4]byte, timeout time.Duration) bool {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Println("failed to create raw socket ", err)
		return false
	}
	defer syscall.Close(fd)
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
