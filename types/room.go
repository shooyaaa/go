package types

import (
	"errors"
	"time"
)

type Room struct {
	members   map[*Session]*Player
	MaxMember int16
	ticker    *time.Ticker
	MsgChan   chan Op
	GameType  Game
	Interval  uint16
	FrameTime int64
	msgBuffer []Op
}

func (r *Room) Init() {
	if r.Interval == 0 {
		r.Interval = 50
	}
	r.ticker = time.NewTicker(time.Duration(r.Interval) * time.Millisecond)
	if r.MaxMember == 0 {
		r.MaxMember = 2000
	}
	r.MsgChan = make(chan Op, 100)
	r.members = make(map[*Session]*Player)
	go r.Tick()
}

func (r *Room) resetMsgChan() {
	r.msgBuffer = make([]Op, 100)
}

func (r *Room) Add(s *Session) error {
	count := int16(len(r.members))
	if count >= r.MaxMember {
		return errors.New("Room fulled")
	}
	r.members[s] = r.GameType.GameData()
	return nil
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

func (r *Room) GetMsgBuffer() []Op {
	return r.msgBuffer
}

func (r *Room) MemberCount() int {
	return len(r.members)
}

func (r *Room) Tick() {
	/*for {
		select {
		case <-r.ticker.C:
			now := time.Now().UnixNano() / 1000000
			r.FrameTime = now - int64(r.Interval)
			r.GameType.HandleOps(r)
			ops := r.GameType.Sync(r)
			for session, _ := range r.members {
				if session.Status == Close {
					delete(r.members, session)
					continue
				}
				_, err := session.Write(ops)
				if err != nil {
					log.Printf("write error %v", err)
				}
			}
			r.resetMsgChan()
		case op := <-r.MsgChan:
			if op.Ts >= r.FrameTime {
				r.msgBuffer = append(r.msgBuffer, op)
			}
		}

	}*/
}
