package manager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shooyaaa/types"
)

var onceSession sync.Once

var instanceSession session

func SessionManager() *session {
	onceSession.Do(func() {
		instanceSession = session{}
		instanceSession.Init()
	})
	return &instanceSession
}

type session struct {
	list     map[types.ID]*types.Session
	WaitChan chan *types.Session
	OpChan   chan types.Op
}

func (s *session) Init() {
	s.list = make(map[types.ID]*types.Session)
	s.WaitChan = make(chan *types.Session, 1000)
	s.OpChan = make(chan types.Op, 1000)
}

func (s *session) Work() {
	fmt.Printf("Start work")
	go s.Accept()
	go s.HandleOp()
}

func (s *session) Accept() {
	for {
		select {
		case session := <-s.WaitChan:
			s.list[session.Id] = session
			session.SetPipe(s.OpChan)
			fmt.Printf("New session %d", session.Id)
		}
	}
}

func (s *session) Get(id types.ID) (*types.Session, error) {
	sessionItem, ok := s.list[id]
	if !ok {
		return sessionItem, errors.New("Invalid Session id %d")
	}
	return sessionItem, nil
}

func (s *session) Push(id types.ID, ops []types.Op) error {
	session, ok := s.list[id]
	if !ok {
		return errors.New("Invalid Session id %d")
	}
	for op := range ops {
		bytes, _ := session.WriteBuffer.Encode(op)
		session.WriteBuffer.Append(bytes)
	}
	return nil
}

func (s *session) HandleOp() {
	for {
		select {
		case op := <-s.OpChan:
			switch op.Type {
			case 1:
				fmt.Println("create room request")
				id, room := RoomManager().Add()
				room.Add(op.Id)
				op.Id.JoinRoom(id, room.MsgChan)
			}
		}
	}
}
