package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/netip"
)

var (
	// ErrInvalidHardwareAddr is returned when one or more invalid hardware
	// addresses are passed to NewPacket.
	ErrInvalidHardwareAddr = errors.New("invalid hardware address")

	// ErrInvalidIP is returned when one or more invalid IPv4 addresses are
	// passed to NewPacket.
	ErrInvalidIP = errors.New("invalid IPv4 address")

	// errInvalidARPPacket is returned when an ethernet frame does not
	// indicate that an ARP packet is contained in its payload.
	errInvalidARPPacket = errors.New("invalid ARP packet")

	Broadcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	ErrInvalidVLAN = errors.New("invalid VLAN")

	ErrInvalidFCS = errors.New("invalid frame check sequence")
)

type EtherType uint16

const (
	EtherTypeIPv4 EtherType = 0x0800
	EtherTypeARP  EtherType = 0x0806
	EtherTypeIPv6 EtherType = 0x86DD

	// EtherTypeVLAN and EtherTypeServiceVLAN are used as 802.1Q Tag Protocol
	// Identifiers (TPIDs).
	EtherTypeVLAN        EtherType = 0x8100
	EtherTypeServiceVLAN EtherType = 0x88a8
	// VLANNone is a special VLAN ID which indicates that no VLAN is being
	// used in a Frame.  In this case, the VLAN's other fields may be used
	// to indicate a Frame's priority.
	VLANNone = 0x000

	// VLANMax is a reserved VLAN ID which may indicate a wildcard in some
	// management systems, but may not be configured or transmitted in a
	// VLAN tag.
	VLANMax = 0xfff

	minPayload = 46

	PriorityBackground           Priority = 1
	PriorityBestEffort           Priority = 0
	PriorityExcellentEffort      Priority = 2
	PriorityCriticalApplications Priority = 3
	PriorityVideo                Priority = 4
	PriorityVoice                Priority = 5
	PriorityInternetworkControl  Priority = 6
	PriorityNetworkControl       Priority = 7
)

//go:generate stringer -output=string.go -type=Operation

// An Operation is an ARP operation, such as request or reply.
type Operation uint16

// Operation constants which indicate an ARP request or reply.
const (
	OperationRequest Operation = 1
	OperationReply   Operation = 2
)

type Priority uint8

type VLAN struct {
	// Priority specifies a IEEE P802.1p priority level.  Priority can be any
	// value from 0 to 7.
	Priority Priority

	// DropEligible indicates if a Frame is eligible to be dropped in the
	// presence of network congestion.
	DropEligible bool

	// ID specifies the VLAN ID for a Frame.  ID can be any value from 0 to
	// 4094 (0x000 to 0xffe), allowing up to 4094 VLANs.
	//
	// If ID is 0 (0x000, VLANNone), no VLAN is specified, and the other fields
	// simply indicate a Frame's priority.
	ID uint16
}

// MarshalBinary allocates a byte slice and marshals a VLAN into binary form.
func (v *VLAN) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	_, err := v.read(b)
	return b, err
}

// read reads data from a VLAN into b.  read is used to marshal a VLAN into
// binary form, but does not allocate on its own.
func (v *VLAN) read(b []byte) (int, error) {
	// Check for VLAN priority in valid range
	if v.Priority > PriorityNetworkControl {
		return 0, ErrInvalidVLAN
	}

	// Check for VLAN ID in valid range
	if v.ID >= VLANMax {
		return 0, ErrInvalidVLAN
	}

	// 3 bits: priority
	ub := uint16(v.Priority) << 13

	// 1 bit: drop eligible
	var drop uint16
	if v.DropEligible {
		drop = 1
	}
	ub |= drop << 12

	// 12 bits: VLAN ID
	ub |= v.ID

	binary.BigEndian.PutUint16(b, ub)
	return 2, nil
}

// UnmarshalBinary unmarshals a byte slice into a VLAN.
func (v *VLAN) UnmarshalBinary(b []byte) error {
	// VLAN tag is always 2 bytes
	if len(b) != 2 {
		return io.ErrUnexpectedEOF
	}

	//  3 bits: priority
	//  1 bit : drop eligible
	// 12 bits: VLAN ID
	ub := binary.BigEndian.Uint16(b[0:2])
	v.Priority = Priority(uint8(ub >> 13))
	v.DropEligible = ub&0x1000 != 0
	v.ID = ub & 0x0fff

	// Check for VLAN ID in valid range
	if v.ID >= VLANMax {
		return ErrInvalidVLAN
	}

	return nil
}

type Frame struct {
	// Destination specifies the destination hardware address for this Frame.
	//
	// If this address is set to Broadcast, the Frame will be sent to every
	// device on a given LAN segment.
	Destination net.HardwareAddr

	// Source specifies the source hardware address for this Frame.
	//
	// Typically, this is the hardware address of the network interface used to
	// send this Frame.
	Source net.HardwareAddr

	// ServiceVLAN specifies an optional 802.1Q service VLAN tag, for use with
	// 802.1ad double tagging, or "Q-in-Q". If ServiceVLAN is not nil, VLAN must
	// not be nil as well.
	//
	// Most users should leave this field set to nil and use VLAN instead.
	ServiceVLAN *VLAN

	// VLAN specifies an optional 802.1Q customer VLAN tag, which may or may
	// not be present in a Frame.  It is important to note that the operating
	// system may automatically strip VLAN tags before they can be parsed.
	VLAN *VLAN

	// EtherType is a value used to identify an upper layer protocol
	// encapsulated in this Frame.
	EtherType EtherType

	// Payload is a variable length data payload encapsulated by this Frame.
	Payload []byte
}

// MarshalBinary allocates a byte slice and marshals a Frame into binary form.
func (f *Frame) MarshalBinary() ([]byte, error) {
	b := make([]byte, f.length())
	_, err := f.read(b)
	return b, err
}

// MarshalFCS allocates a byte slice, marshals a Frame into binary form, and
// finally calculates and places a 4-byte IEEE CRC32 frame check sequence at
// the end of the slice.
//
// Most users should use MarshalBinary instead.  MarshalFCS is provided as a
// convenience for rare occasions when the operating system cannot
// automatically generate a frame check sequence for an Ethernet frame.
func (f *Frame) MarshalFCS() ([]byte, error) {
	// Frame length with 4 extra bytes for frame check sequence
	b := make([]byte, f.length()+4)
	if _, err := f.read(b); err != nil {
		return nil, err
	}

	// Compute IEEE CRC32 checksum of frame bytes and place it directly
	// in the last four bytes of the slice
	binary.BigEndian.PutUint32(b[len(b)-4:], crc32.ChecksumIEEE(b[0:len(b)-4]))
	return b, nil
}

// read reads data from a Frame into b.  read is used to marshal a Frame
// into binary form, but does not allocate on its own.
func (f *Frame) read(b []byte) (int, error) {
	// S-VLAN must also have accompanying C-VLAN.
	if f.ServiceVLAN != nil && f.VLAN == nil {
		return 0, ErrInvalidVLAN
	}

	copy(b[0:6], f.Destination)
	copy(b[6:12], f.Source)

	// Marshal each non-nil VLAN tag into bytes, inserting the appropriate
	// EtherType/TPID before each, so devices know that one or more VLANs
	// are present.
	vlans := []struct {
		vlan *VLAN
		tpid EtherType
	}{
		{vlan: f.ServiceVLAN, tpid: EtherTypeServiceVLAN},
		{vlan: f.VLAN, tpid: EtherTypeVLAN},
	}

	n := 12
	for _, vt := range vlans {
		if vt.vlan == nil {
			continue
		}

		// Add VLAN EtherType and VLAN bytes.
		binary.BigEndian.PutUint16(b[n:n+2], uint16(vt.tpid))
		if _, err := vt.vlan.read(b[n+2 : n+4]); err != nil {
			return 0, err
		}
		n += 4
	}

	// Marshal actual EtherType after any VLANs, copy payload into
	// output bytes.
	binary.BigEndian.PutUint16(b[n:n+2], uint16(f.EtherType))
	copy(b[n+2:], f.Payload)

	return len(b), nil
}

// UnmarshalBinary unmarshals a byte slice into a Frame.
func (f *Frame) UnmarshalBinary(b []byte) error {
	// Verify that both hardware addresses and a single EtherType are present
	if len(b) < 14 {
		return io.ErrUnexpectedEOF
	}

	// Track offset in packet for reading data
	n := 14

	// Continue looping and parsing VLAN tags until no more VLAN EtherType
	// values are detected
	et := EtherType(binary.BigEndian.Uint16(b[n-2 : n]))
	switch et {
	case EtherTypeServiceVLAN, EtherTypeVLAN:
		// VLAN type is hinted for further parsing.  An index is returned which
		// indicates how many bytes were consumed by VLAN tags.
		nn, err := f.unmarshalVLANs(et, b[n:])
		if err != nil {
			return err
		}

		n += nn
	default:
		// No VLANs detected.
		f.EtherType = et
	}

	// Allocate single byte slice to store destination and source hardware
	// addresses, and payload
	bb := make([]byte, 6+6+len(b[n:]))
	copy(bb[0:6], b[0:6])
	f.Destination = bb[0:6]
	copy(bb[6:12], b[6:12])
	f.Source = bb[6:12]

	// There used to be a minimum payload length restriction here, but as
	// long as two hardware addresses and an EtherType are present, it
	// doesn't really matter what is contained in the payload.  We will
	// follow the "robustness principle".
	copy(bb[12:], b[n:])
	f.Payload = bb[12:]

	return nil
}

// UnmarshalFCS computes the IEEE CRC32 frame check sequence of a Frame,
// verifies it against the checksum present in the byte slice, and finally,
// unmarshals a byte slice into a Frame.
//
// Most users should use UnmarshalBinary instead.  UnmarshalFCS is provided as
// a convenience for rare occasions when the operating system cannot
// automatically verify a frame check sequence for an Ethernet frame.
func (f *Frame) UnmarshalFCS(b []byte) error {
	// Must contain enough data for FCS, to avoid panics
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}

	// Verify checksum in slice versus newly computed checksum
	want := binary.BigEndian.Uint32(b[len(b)-4:])
	got := crc32.ChecksumIEEE(b[0 : len(b)-4])
	if want != got {
		return ErrInvalidFCS
	}

	return f.UnmarshalBinary(b[0 : len(b)-4])
}

// length calculates the number of bytes required to store a Frame.
func (f *Frame) length() int {
	// If payload is less than the required minimum length, we zero-pad up to
	// the required minimum length
	pl := len(f.Payload)
	if pl < minPayload {
		pl = minPayload
	}

	// Add additional length if VLAN tags are needed.
	var vlanLen int
	switch {
	case f.ServiceVLAN != nil && f.VLAN != nil:
		vlanLen = 8
	case f.VLAN != nil:
		vlanLen = 4
	}

	// 6 bytes: destination hardware address
	// 6 bytes: source hardware address
	// N bytes: VLAN tags (if present)
	// 2 bytes: EtherType
	// N bytes: payload length (may be padded)
	return 6 + 6 + vlanLen + 2 + pl
}

// unmarshalVLANs unmarshals S/C-VLAN tags.  It is assumed that tpid
// is a valid S/C-VLAN TPID.
func (f *Frame) unmarshalVLANs(tpid EtherType, b []byte) (int, error) {
	// 4 or more bytes must remain for valid S/C-VLAN tag and EtherType.
	if len(b) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	// Track how many bytes are consumed by VLAN tags.
	var n int

	switch tpid {
	case EtherTypeServiceVLAN:
		vlan := new(VLAN)
		if err := vlan.UnmarshalBinary(b[n : n+2]); err != nil {
			return 0, err
		}
		f.ServiceVLAN = vlan

		// Assume that a C-VLAN immediately trails an S-VLAN.
		if EtherType(binary.BigEndian.Uint16(b[n+2:n+4])) != EtherTypeVLAN {
			return 0, ErrInvalidVLAN
		}

		// 4 or more bytes must remain for valid C-VLAN tag and EtherType.
		n += 4
		if len(b[n:]) < 4 {
			return 0, io.ErrUnexpectedEOF
		}

		// Continue to parse the C-VLAN.
		fallthrough
	case EtherTypeVLAN:
		vlan := new(VLAN)
		if err := vlan.UnmarshalBinary(b[n : n+2]); err != nil {
			return 0, err
		}

		f.VLAN = vlan
		f.EtherType = EtherType(binary.BigEndian.Uint16(b[n+2 : n+4]))
		n += 4
	default:
		panic(fmt.Sprintf("unknown VLAN TPID: %04x", tpid))
	}

	return n, nil
}

// A Packet is a raw ARP packet, as described in RFC 826.
type Packet struct {
	// HardwareType specifies an IANA-assigned hardware type, as described
	// in RFC 826.
	HardwareType uint16

	// ProtocolType specifies the internetwork protocol for which the ARP
	// request is intended.  Typically, this is the IPv4 EtherType.
	ProtocolType uint16

	// HardwareAddrLength specifies the length of the sender and target
	// hardware addresses included in a Packet.
	HardwareAddrLength uint8

	// IPLength specifies the length of the sender and target IPv4 addresses
	// included in a Packet.
	IPLength uint8

	// Operation specifies the ARP operation being performed, such as request
	// or reply.
	Operation Operation

	// SenderHardwareAddr specifies the hardware address of the sender of this
	// Packet.
	SenderHardwareAddr net.HardwareAddr

	// SenderIP specifies the IPv4 address of the sender of this Packet.
	SenderIP netip.Addr

	// TargetHardwareAddr specifies the hardware address of the target of this
	// Packet.
	TargetHardwareAddr net.HardwareAddr

	// TargetIP specifies the IPv4 address of the target of this Packet.
	TargetIP netip.Addr
}

// NewPacket creates a new Packet from an input Operation and hardware/IPv4
// address values for both a sender and target.
//
// If either hardware address is less than 6 bytes in length, or there is a
// length mismatch between the two, ErrInvalidHardwareAddr is returned.
//
// If either IP address is not an IPv4 address, or there is a length mismatch
// between the two, ErrInvalidIP is returned.
func NewPacket(op Operation, srcHW net.HardwareAddr, srcIP netip.Addr, dstHW net.HardwareAddr, dstIP netip.Addr) (*Packet, error) {
	// Validate hardware addresses for minimum length, and matching length
	if len(srcHW) < 6 {
		return nil, ErrInvalidHardwareAddr
	}
	if len(dstHW) < 6 {
		return nil, ErrInvalidHardwareAddr
	}
	if !bytes.Equal(Broadcast, dstHW) && len(srcHW) != len(dstHW) {
		return nil, ErrInvalidHardwareAddr
	}

	// Validate IP addresses to ensure they are IPv4 addresses, and
	// correct length
	var invalidIP netip.Addr
	if !srcIP.IsValid() || !srcIP.Is4() {
		return nil, ErrInvalidIP
	}
	if !dstIP.Is4() || dstIP == invalidIP {
		return nil, ErrInvalidIP
	}

	return &Packet{
		// There is no Go-native way to detect hardware type of a network
		// interface, so default to 1 (ethernet 10Mb) for now
		HardwareType: 1,

		// Default to EtherType for IPv4
		ProtocolType: uint16(EtherTypeARP),

		// Populate other fields using input data
		HardwareAddrLength: uint8(len(srcHW)),
		IPLength:           uint8(4),
		Operation:          op,
		SenderHardwareAddr: srcHW,
		SenderIP:           srcIP,
		TargetHardwareAddr: dstHW,
		TargetIP:           dstIP,
	}, nil
}

func (p *Packet) Marshal() ([]byte, error) {

	return nil, nil
}

// MarshalBinary allocates a byte slice containing the data from a Packet.
//
// MarshalBinary never returns an error.
func (p *Packet) MarshalBinary() ([]byte, error) {
	// 2 bytes: hardware type
	// 2 bytes: protocol type
	// 1 byte : hardware address length
	// 1 byte : protocol length
	// 2 bytes: operation
	// N bytes: source hardware address
	// N bytes: source protocol address
	// N bytes: target hardware address
	// N bytes: target protocol address

	// Though an IPv4 address should always 4 bytes, go-fuzz
	// very quickly created several crasher scenarios which
	// indicated that these values can lie.
	b := make([]byte, 2+2+1+1+2+(p.IPLength*2)+(p.HardwareAddrLength*2))

	// Marshal fixed length data

	binary.BigEndian.PutUint16(b[0:2], p.HardwareType)
	binary.BigEndian.PutUint16(b[2:4], p.ProtocolType)

	b[4] = p.HardwareAddrLength
	b[5] = p.IPLength

	binary.BigEndian.PutUint16(b[6:8], uint16(p.Operation))

	// Marshal variable length data at correct offset using lengths
	// defined in p

	n := 8
	hal := int(p.HardwareAddrLength)
	pl := int(p.IPLength)

	copy(b[n:n+hal], p.SenderHardwareAddr)
	n += hal

	sender4 := p.SenderIP.As4()
	copy(b[n:n+pl], sender4[:])
	n += pl

	copy(b[n:n+hal], p.TargetHardwareAddr)
	n += hal

	target4 := p.TargetIP.As4()
	copy(b[n:n+pl], target4[:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a Packet.
func (p *Packet) UnmarshalBinary(b []byte) error {
	// Must have enough room to retrieve hardware address and IP lengths
	if len(b) < 8 {
		return io.ErrUnexpectedEOF
	}

	// Retrieve fixed length data

	p.HardwareType = binary.BigEndian.Uint16(b[0:2])
	p.ProtocolType = binary.BigEndian.Uint16(b[2:4])

	p.HardwareAddrLength = b[4]
	p.IPLength = b[5]

	p.Operation = Operation(binary.BigEndian.Uint16(b[6:8]))

	// Unmarshal variable length data at correct offset using lengths
	// defined by ml and il
	//
	// These variables are meant to improve readability of offset calculations
	// for the code below
	n := 8
	ml := int(p.HardwareAddrLength)
	ml2 := ml * 2
	il := int(p.IPLength)
	il2 := il * 2

	// Must have enough room to retrieve both hardware address and IP addresses
	addrl := n + ml2 + il2
	if len(b) < addrl {
		return io.ErrUnexpectedEOF
	}

	// Allocate single byte slice to store address information, which
	// is resliced into fields
	bb := make([]byte, addrl-n)

	// Sender hardware address
	copy(bb[0:ml], b[n:n+ml])
	p.SenderHardwareAddr = bb[0:ml]
	n += ml

	// Sender IP address
	copy(bb[ml:ml+il], b[n:n+il])
	senderIP, ok := netip.AddrFromSlice(bb[ml : ml+il])
	if !ok {
		return errors.New("Invalid Sender IP address")
	}
	p.SenderIP = senderIP
	n += il

	// Target hardware address
	copy(bb[ml+il:ml2+il], b[n:n+ml])
	p.TargetHardwareAddr = bb[ml+il : ml2+il]
	n += ml

	// Target IP address
	copy(bb[ml2+il:ml2+il2], b[n:n+il])
	targetIP, ok := netip.AddrFromSlice(bb[ml2+il : ml2+il2])
	if !ok {
		return errors.New("Invalid Target IP address")
	}
	p.TargetIP = targetIP

	return nil
}

func parsePacket(buf []byte) (*Packet, *Frame, error) {
	f := new(Frame)
	if err := f.UnmarshalBinary(buf); err != nil {
		return nil, nil, err
	}

	// Ignore frames which do not have ARP EtherType
	if f.EtherType != EtherTypeARP {
		return nil, nil, errInvalidARPPacket
	}

	p := new(Packet)
	if err := p.UnmarshalBinary(f.Payload); err != nil {
		return nil, nil, err
	}
	return p, f, nil
}
