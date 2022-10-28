package unimouse

import (
	"fmt"

	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/op"
	"github.com/shooyaaa/core/session"

	. "github.com/shooyaaa/core/types"
)

func Run() {
	tcp := network.Tcp{Id: &Simple{}}
	tcp.Listen(":9994")
	for {
		handler := Handler{}
		session := tcp.Accept()
		session.SetOwner(handler)
		handler.clients[session.Id] = session
	}
}

type Handler struct {
	clients map[int64]*session.Session
}

func (h Handler) OpHandler(op op.Op, s *session.Session) {
	fmt.Println("op comes ", op, " session ", s)
}
func (h Handler) SessionClose(id int64) {
	delete(h.clients, id)
}
