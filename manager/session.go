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

var instanceSession sm

func SessionManager() *sm {
	onceSession.Do(func() {
		instanceSession = sm{}
		instanceSession.Init()
	})
	return &instanceSession
}

type sm struct {
	list     map[int64]*types.Session
	WaitChan chan *types.Session
	OpChan   chan types.Op
}

func (s *sm) Init() {
	s.list = make(map[int64]*types.Session)
	s.WaitChan = make(chan *types.Session, 1000)
	s.OpChan = make(chan types.Op, 1000)
}

func (s *sm) Work() {
	fmt.Printf("Start work")
	go s.Accept()
	go s.HandleOp()
}

func (s *sm) Accept() {
	for {
		select {
		case session := <-s.WaitChan:
			s.list[session.Id] = session
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

func (s *sm) Get(id int64) (*types.Session, error) {
	sessionItem, ok := s.list[id]
	if !ok {
		return sessionItem, errors.New("Invalid Session id %d")
	}
	return sessionItem, nil
}

func (s *sm) Push(id int64, ops []types.Op) error {
	session, ok := s.list[id]
	if !ok {
		return errors.New("Invalid Session id %d")
	}
	session.Write(ops)
	return nil
}

func (s *sm) HandleOp() {
	for {
		select {
		case op := <-s.OpChan:
			session := op.GetId()
			switch op.Type {
			case types.Op_Create_Room:
				fmt.Println("create room request")
				_, room := RoomManager().Add()
				room.GameType = game.Snake{}
				room.Add(session)
			case types.Op_Join_Room:
				d := op.Data["Id"]
				roomId := int64(d)
				room, err := RoomManager().Get(roomId)
				if err != nil {
					log.Printf("Error while join room %v", err)
				}
				room.Add(session)
			case types.Op_Logout:
				id := int64(op.Data["Id"])
				delete(s.list, id)
			default:
				log.Printf("op comes here")
			}
		}
	}
}
