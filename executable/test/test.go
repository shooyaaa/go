package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/shooyaaa/core/network"
)

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
*/
import "C"

func main() {
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
	ip_cstr := C.CString("192.168.50.65")

	packet := C.GoBytes(unsafe.Pointer(C.FillRequestPacketFields(iface_cstr, ip_cstr)), C.int(size))

	// Send the packet
	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ARP
	addr.Ifindex = interf.Index
	addr.Hatype = syscall.ARPHRD_ETHER

	//err = syscall.Sendto(fd, packet, 0, &addr)
	err = network.SendRaw(packet)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Sent packet")
	}

}
