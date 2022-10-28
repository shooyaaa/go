package unimouse

import (
	"fmt"

	. "github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/op"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/log"
)

func Connect() {
	tcp := TcpConn{}
	session, err := tcp.Dial("localhost", 9994)
	if err != nil {
		log.Fatal("error occurs while connect to server %v", err)
	}

	handler := Handler{}
	session.SetOwner(handler)
}

type ClientHandler struct {
}

func (c ClientHandler) OpHandler(op op.Op, s *session.Session) {
	fmt.Println("op comes ", op, " session ", s)
}
func (c ClientHandler) SessionClose(id int64) {
	log.InfoF("connection closed %v", id)
}
