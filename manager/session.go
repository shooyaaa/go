package manager

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/shooyaaa/game"

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
			session.SetPipe(&s.OpChan)
			fmt.Printf("New session %d", session.Id)
			ops := make([]types.Op, 1)
			data := make(map[string]float64)
			data["id"] = float64(session.Id)
			ops[0] = types.Op{
				Type: types.Op_Login,
				Data: data,
			}
			session.Write(ops)
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
	bytes, _ := session.WriteBuffer.Encode(ops)
	session.WriteBuffer.Append(bytes)
	return nil
}

func (s *session) HandleOp() {
	for {
		select {
		case op := <-s.OpChan:
			s := op.GetId()
			switch op.Type {
			case 1:
				fmt.Println("create room request")
				id, room := RoomManager().Add()
				room.GameType = game.Snake{}
				room.Add(s)
				s.JoinRoom(id, &room.MsgChan)
			case 2:
				d := op.Data["Id"]
				roomId := types.ID(d)
				room, err := RoomManager().Get(roomId)
				if err != nil {
					log.Printf("Error while join room %v", err)
				}
				room.Add(s)
				s.JoinRoom(roomId, &room.MsgChan)
			default:
				log.Printf("op comes here")
			}
		}
	}
}
