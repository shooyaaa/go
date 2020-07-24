package types

import (
	"errors"
	"log"
	"time"
)

type Room struct {
	members   map[*Session]Player
	MaxMember int16
	ticker    *time.Ticker
	MsgChan   chan Op
	GameType  Game
	Interval  uint16
	frameTime int64
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
	r.members = make(map[*Session]Player)
	go r.Tick()
}

func (r *Room) Add(id *Session) error {
	count := int16(len(r.members))
	if count >= r.MaxMember {
		return errors.New("Room fulled")
	}
	r.members[id] = Player{} //r.GameType.GameData()
	return nil
}

func (r *Room) Leave(id *Session) error {
	_, ok := r.members[id]
	if !ok {
		return errors.New("Player not found in Room ")
	}
	delete(r.members, id)
	return nil
}

func (r *Room) Tick() {
	for {
		select {
		case <-r.ticker.C:
			now := time.Now().UnixNano() / 1000
			r.frameTime = now - int64(r.Interval)
			//r.GameType.HandleOps(r.msgBuffer)
			//ops := r.GameType.Sync()
			ops := make([]Op, 1)
			ops = append(ops, Op{Type: 1})
			for op := range ops {
				for session, _ := range r.members {
					//bytes, _ := session.WriteBuffer.Encode(op)
					//session.WriteBuffer.Append(bytes)
					err := session.Write(op)
					if err != nil {
						log.Printf("write error %v", err)
					}
				}
			}
		case op := <-r.MsgChan:
			if op.Ts >= r.frameTime {
				r.msgBuffer = append(r.msgBuffer, op)
			}
		}

	}
}
