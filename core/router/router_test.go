package router

import (
	"fmt"
	"github.com/shooyaaa/core"
	"testing"
)

func TestTcpRouter_Deliver(t *testing.T) {
	header := Header{
		seq: 0,
		ack: 0,
		len: 0,
	}
	p := Package{
		Header: header,
		body:   []byte("hello world"),
	}
	r, err := LookUp("0:xiaocui")
	if err.Error() != core.NOT_FOUND {
		t.Fatalf("err not correct: %v", err.Error())
	}
	if r != nil {
		t.Fatalf("router shoulde be nil %v", r)
	}
	host := "127.0.0.1"
	port := "8087"
	AddTcpAddress("xiaocui", host+":"+port)
	r, err = LookUp("0:xiaocui")
	tcp, ok := r.(*TcpRouter)
	if err != nil || ok == false {
		t.Fatalf("router not found %v", err)
	}
	tcp.Listen("127.0.0.1:8087")
	if r.ToString() != fmt.Sprintf("%v:%v:%v", TCP_ROUTER, host, port) {
		t.Fatalf("the address found is not correct %v", r.ToString())
	}
	err = tcp.Forward(&p)
	if err != nil {
		t.Log("error should be nil")
	}
}
