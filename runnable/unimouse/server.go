package unimouse

import (
	"fmt"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/log"

	. "github.com/shooyaaa/core/types"
)

func Run() {
	tcp := network.Tcp{Id: &Simple{}}
	tcp.Listen(":9994")
	session.SetCodec(codec.NewCodec[codec.Op](codec.JSON_CODEC))
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
	var lastTime time.Time = time.Now()
	hook.Register(hook.MouseMove, []string{}, func(e hook.Event) {
		o := codec.MakeOp(codec.Op_MouseEvent, map[string]interface{}{
			"X": e.X, "Y": e.Y,
		})
		diff := e.When.Sub(lastTime)
		if diff > 5*time.Microsecond {
			log.DebugF("mouse pos x: %v, y: %v, diff: %v", e.X, e.Y, diff)
			h.Manager.Broadcast(o)
			lastTime = e.When
			hook.StopEvent()
		}
	})
	hook.Register(hook.KeyUp, []string{}, func(e hook.Event) {
		o := codec.MakeOp(codec.Op_KeyEvent, map[string]interface{}{
			"RawCode": e.Rawcode, "Keychar": e.Keychar,
		})
		h.Manager.Broadcast(o)
		hook.StopEvent()
	})
	s := hook.Start()
	<-hook.Process(s)
}

type Handler struct {
	session.Manager
}

func (h Handler) OpHandler(op codec.Op, s *session.Session) {
	fmt.Println("server op comes ", op, " session ", s)
}
func (h Handler) SessionClose(id int64) {
	h.Manager.RemoveId(id)
}
