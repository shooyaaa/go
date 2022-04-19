package types

import (
	"errors"
	"github.com/shooyaaa/log"
	"time"
)

type Room struct {
	members   map[*Session]*Player
	MaxMember int16
	ticker    *time.Ticker
	MsgChan   chan OpWithSession
	GameType  Game
	Interval  uint16
	FrameTime int64
	msgBuffer []OpWithSession
}

func (r *Room) Init() {
	if r.Interval == 0 {
		r.Interval = 50
	}
	r.ticker = time.NewTicker(time.Duration(r.Interval) * time.Millisecond)
	if r.MaxMember == 0 {
		r.MaxMember = 2000
	}
	r.MsgChan = make(chan OpWithSession, 100)
	r.members = make(map[*Session]*Player)
	go r.Tick()
}

func (r *Room) resetMsgChan() {
	r.msgBuffer = make([]OpWithSession, 100)
}

func (r *Room) Add(s *Session) error {
	count := int16(len(r.members))
	if count >= r.MaxMember {
		return errors.New("Room fulled")
	}
	r.members[s] = r.GameType.GameData()
	s.SetOwner(r)
	return nil
}
func (r *Room) OpHandler(op Op, session *Session) {
	opSession := OpWithSession{
		Op:      op,
		session: session,
	}
	r.MsgChan <- opSession
}

func (r *Room) OpHandler1(op Op, session *Session) {
	switch op.Type {
	case Op_Logout:
		delete(r.members, session)
	case Op_Sync_Data:
		gameData := r.members[session]
		x, ok := op.Data["x"]
		if ok {
			gameData.X = x
		}
		y, ok := op.Data["y"]
		if ok {
			gameData.Y = y
		}
		log.DebugF("Player %v moved to x: %v, y : %v", session.Id, gameData.X, gameData.Y)
	default:
		log.WarnF("unhandled op in room %v", op.Type)
	}
}

func (r *Room) SessionClose(id int64) {

}

func (r *Room) Leave(id *Session) error {
	_, err := r.GetMember(id)
	if err != nil {
		return err
	}
	delete(r.members, id)
	return nil
}

func (r *Room) GetMember(id *Session) (*Player, error) {
	data, ok := r.members[id]
	if !ok {
		return nil, errors.New("Player not found in Room ")
	}
	return data, nil
}

func (r *Room) AllMembers() map[*Session]*Player {
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
