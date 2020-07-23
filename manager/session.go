package manager

import (
	"github.com/shooyaaa/types"
	"fmt"
)

type Session struct {
	list 		map[types.ID]types.Session
	WaitChan 	chan types.Session
	OpChan		chan types.Op
}


func (s *Session) Init () {
	s.list = make(map[types.ID]types.Session)
}

func (s *Session) Work () {
	fmt.Printf("Start work");
	go s.Accept()
	go s.HandleOp()
}

func (s *Session) Accept () {
	for {
		select {
		case session := <- s.WaitChan:
			s.list[session.Id] = session
			session.OpPipe = s.OpChan
			fmt.Printf("New session %d", session.Id);
		}
	}
}

func (s *Session) HandleOp() {
	for {
		select {
		case op := <- s.OpChan:
			switch op.Type {
			case 1:
				fmt.Println("create room request")	
			}
		}
	}	
}