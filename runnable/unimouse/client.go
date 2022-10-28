package unimouse

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-vgo/robotgo"
	. "github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/log"
)

func Connect() {
	tcp := TcpConn{}
	s, err := tcp.Dial("10.1.36.39", 9994)
	if err != nil {
		log.Fatal("error occurs while connect to server %v", err)
		return
	}
	session.SetCodec(&session.Json{})
	handler := ClientHandler{}
	s.SetOwner(handler)
	ch := make(chan os.Signal, 1)
	log.Info(http.ListenAndServe("localhost:6060", nil).Error())
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}

type ClientHandler struct {
}

func (c ClientHandler) OpHandler(op session.Op, s *session.Session) {
	fmt.Println("op comes ", op.Type, " session ", s.Id)
	switch op.Type {
	case session.Op_KeyEvent:
		log.DebugF("print key %v ", op.Data["KeyChar"])
	case session.Op_MouseEvent:
		x := op.Data["X"]
		y := op.Data["Y"]
		robotgo.Move(int(x.(float64)), int(y.(float64)))
	}
}
func (c ClientHandler) SessionClose(id int64) {
	log.InfoF("connection closed %v", id)
}
