package types

import (
	"errors"
	"time"

	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/log"
)

type Room struct {
	members   map[*session.Session]*Player
	MaxMember int16
	ticker    *time.Ticker
	MsgChan   chan codec.Op
	GameType  Game
	Interval  uint16
	FrameTime int64
	msgBuffer []codec.Op
}

func (r *Room) Init() {
	if r.Interval == 0 {
		r.Interval = 50
	}
	r.ticker = time.NewTicker(time.Duration(r.Interval) * time.Millisecond)
	if r.MaxMember == 0 {
		r.MaxMember = 2000
	}
	r.MsgChan = make(chan codec.Op, 100)
	r.members = make(map[*session.Session]*Player)
	go r.Tick()
}

func (r *Room) resetMsgChan() {
	r.msgBuffer = make([]codec.Op, 100)
}

func (r *Room) Add(s *session.Session) error {
	count := int16(len(r.members))
	if count >= r.MaxMember {
		return errors.New("Room fulled")
	}
	r.members[s] = nil
	s.SetOwner(r)
	return nil
}
func (r *Room) OpHandler(op codec.Op, session *session.Session) {
	r.MsgChan <- op
}

func (r *Room) OpHandler1(op1 codec.Op, s *session.Session) {
	switch op1.Type {
	case codec.Op_Logout:
		delete(r.members, s)
	case codec.Op_Sync_Data:
		gameData := r.members[s]
		x, ok := op1.Data["x"]
		if ok {
			gameData.X = x.(float64)
		}
		y, ok := op1.Data["y"]
		if ok {
			gameData.Y = y.(float64)
		}
		log.DebugF("Player %v moved to x: %v, y : %v", s.Id, gameData.X, gameData.Y)
	default:
		log.WarnF("unhandled op in room %v", op1.Type)
	}
}

func (r *Room) SessionClose(id int64) {

}

func (r *Room) Leave(id *session.Session) error {
	_, err := r.GetMember(id)
	if err != nil {
		return err
	}
	delete(r.members, id)
	return nil
}

func (r *Room) GetMember(id *session.Session) (*Player, error) {
	data, ok := r.members[id]
	if !ok {
		return nil, errors.New("Player not found in Room ")
	}
	return data, nil
}

func (r *Room) AllMembers() map[*session.Session]*Player {
	return r.members
}

func (r *Room) MemberCount() int {
	return len(r.members)
}

func (r *Room) Tick() {
	for {
		select {
		case <-r.ticker.C:
			now := time.Now().UnixNano() / 1000000
			r.FrameTime = now - int64(r.Interval)
			r.GameType.Play(r.msgBuffer)
			r.resetMsgChan()
		case op := <-r.MsgChan:
			if op.Ts >= r.FrameTime {
				r.msgBuffer = append(r.msgBuffer, op)
			}
		}

	}
}
