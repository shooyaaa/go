package unimouse

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/shooyaaa/core/codec"
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
	session.SetCodec(&codec.Json{})
	handler := Handler{}
	s.SetOwner(handler)
	ch := make(chan os.Signal, 1)
	log.Info(http.ListenAndServe("localhost:6060", nil).Error())
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}

type ClientHandler struct {
}

func (c ClientHandler) OpHandler(op session.Op, s *session.Session) {
	fmt.Println("op comes ", op, " session ", s)
}
func (c ClientHandler) SessionClose(id int64) {
	log.InfoF("connection closed %v", id)
}
