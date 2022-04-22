package tests

import (
	. "github.com/shooyaaa/core/codec"
	. "github.com/shooyaaa/core/op"
	"net"
	"testing"
)

func makeOp(t uint8) Op {
	o := Op{
		Type: t,
		Ts:   0,
		Data: map[string]float64{},
	}
	return o
}

func CreateRoom() Op {
	return makeOp(1)
}
func JoinRoom() Op {
	return makeOp(2)
}

func SyncData() Op {
	return makeOp(3)
}

func Login() Op {
	return makeOp(4)
}

func Logout() Op {
	return makeOp(5)
}

func TestServer(t *testing.T) {
	c, err := net.Dial("tcp", "127.0.0.1:3352")
	if err != nil {
		t.Error("error while connect to server")
	}
	codec := Json{}
	sd := SyncData()
	sd.Data = map[string]float64{"X": 1.4, "Y": 1.5}
	b, err := codec.Encode([]Op{CreateRoom(), sd})
	c.Write(b)
}
