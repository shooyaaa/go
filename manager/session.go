package manager

import (
	"errors"
	"github.com/shooyaaa/core/op"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/game"
	"github.com/shooyaaa/log"
	"sync"
)

var onceSession sync.Once

var SessionManager sm

func init() {
	onceSession.Do(func() {
		SessionManager = sm{}
		SessionManager.Init()
	})
}

type sm struct {
	list     map[int64]*session.Session
	WaitChan chan *session.Session
	OpChan   chan op.Op
}

func (s *sm) Init() {
	s.list = make(map[int64]*session.Session)
	s.WaitChan = make(chan *session.Session, 1000)
	s.OpChan = make(chan op.Op, 1000)
}

func (s *sm) Work() {
	log.Debug("session manager start Work")
	go s.Accept()
}

func (s *sm) Accept() {
	for {
		select {
		case session := <-s.WaitChan:
			s.list[session.Id] = session
			session.SetOwner(s)
			log.DebugF("New session %v", session.Id)
			ops := make([]op.Op, 1)
			data := make(map[string]float64)
			data["id"] = float64(session.Id)
			ops[0] = op.Op{
				Type: op.Op_Login,
				Data: data,
			}
			session.Write(ops)
		}
	}
}

func (s *sm) Get(id int64) (*session.Session, error) {
	sessionItem, ok := s.list[id]
	if !ok {
		return sessionItem, errors.New("Invalid Session id %d")
	}
	return sessionItem, nil
}

func (s *sm) Push(id int64, ops []op.Op) error {
	session, ok := s.list[id]
	if !ok {
		return errors.New("Invalid Session id %d")
	}
	session.Write(ops)
	return nil
}

func (s *sm) OpHandler(op op.Op, session *session.Session) {
	switch op.Type {
	case op.Op_Create_Room:
		log.Debug("create room request")
		_, room := RoomManager().Add()
		room.GameType = game.Snake{}
		room.Add(session)
	case op.Op_Join_Room:
		d := op.Data["Id"]
		roomId := int64(d)
		room, err := RoomManager().Get(roomId)
		if err != nil {
			log.DebugF("Error while join room %v", err)
		}
		room.Add(session)
	case op.Op_Logout:
		id := int64(op.Data["Id"])
		delete(s.list, id)
	default:
		log.Debug("session manager op comes here")
	}
}
func (s *sm) SessionClose(id int64) {
	log.DebugF("session close in session manager: %v", id)
	delete(s.list, id)
}
