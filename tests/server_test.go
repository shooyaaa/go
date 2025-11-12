package tests

import (
	"net"
	"testing"

	"github.com/shooyaaa/core/codec"
)

func makeOp(t codec.OpType) codec.Op {
	o := codec.Op{
		Type: t,
		Ts:   0,
		Data: map[string]interface{}{},
	}
	return o
}

func CreateRoom() codec.Op {
	return makeOp(1)
}
func JoinRoom() codec.Op {
	return makeOp(2)
}

func SyncData() codec.Op {
	return makeOp(3)
}

func Login() codec.Op {
	return makeOp(4)
}

func Logout() codec.Op {
	return makeOp(5)
}

func TestServer(t *testing.T) {
	c, err := net.Dial("tcp", "127.0.0.1:3352")
	if err != nil {
		t.Error("error while connect to server")
	}
	codecInstance := codec.NewCodec[codec.Op](codec.JSON_CODEC)
	sd := SyncData()
	sd.Data = map[string]interface{}{"X": 1.4, "Y": 1.5}
	b, err := codecInstance.Encode(codec.Op{Type: codec.Op_Sync_Data, Data: map[string]interface{}{"X": 1.4, "Y": 1.5}})
	if err != nil {
		t.Error("error while encode op")
	}
	c.Write(b)
}
