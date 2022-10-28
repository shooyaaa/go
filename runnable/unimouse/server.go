package unimouse

import (
	"fmt"

	hook "github.com/robotn/gohook"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/session"

	. "github.com/shooyaaa/core/types"
)

func Run() {
	tcp := network.Tcp{Id: &Simple{}}
	tcp.Listen(":9994")
	session.SetCodec(&session.Json{})
	handler := Handler{}
	handler.Manager.Init()
	go startHook(handler)
	for {
		s := tcp.Accept()
		s.SetOwner(handler)
		handler.Manager.Add(*s)
	}
}

func startHook(h Handler) {
	hook.Register(hook.MouseMove, []string{}, func(e hook.Event) {
		o := session.MakeOp(session.Op_MouseEvent, map[string]interface{}{
			"X": e.X, "Y": e.Y,
		})
		h.Manager.Broadcast(o)
	})
	hook.Register(hook.KeyUp, []string{}, func(e hook.Event) {
		o := session.MakeOp(session.Op_KeyEvent, map[string]interface{}{
			"RawCode": e.Rawcode, "Keychar": e.Keychar,
		})
		h.Manager.Broadcast(o)
	})
	s := hook.Start()
	<-hook.Process(s)
}

type Handler struct {
	session.Manager
}

func (h Handler) OpHandler(op session.Op, s *session.Session) {
	fmt.Println("op comes ", op, " session ", s)
}
func (h Handler) SessionClose(id int64) {
	h.Manager.RemoveId(id)
}
