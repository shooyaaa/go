package types

import (
	"time"
	"errors"
)

type Room struct {
	members 	map[ID]Player
	MaxMember 	int16
	ticker   	*time.Ticker
	MsgChan		chan Op
	GameType 	Game
	Interval	uint16
	frameTime 	int64
	msgBuffer   []Op
}

func (r *Room) Init() {
	if r.Interval == 0 {
		r.Interval = 50
	}
	r.ticker = time.NewTicker(time.Duration(r.Interval) * time.Millisecond)
	if r.MaxMember == 0 {
		r.MaxMember = 2000
	}
	r.members = make(map[ID]Player)
}

func (r *Room) Add(id ID) error {
	if int16(len(r.members)) >= r.MaxMember {
		return errors.New("Room fulled")
	}
	r.members[id] = r.GameType.GameData()
	return nil
}

func (r *Room) Leave(id ID) error {
	_, ok := r.members[id]
	if !ok {
		return errors.New("Player not found in Room ")
	}
	delete(r.members, id)
	return nil
}

func (r *Room) Handle(msg interface{}) {

}

func (r *Room) Tick() {
	for  {
		select {
		case <-r.ticker.C:
			now := time.Now().UnixNano() / 1000
			r.frameTime = now - int64(r.Interval)
			r.GameType.HandleOps(r.msgBuffer)
			r.GameType.Sync()
		case op := <-r.MsgChan:
			if op.Ts >= r.frameTime {
				r.msgBuffer = append(r.msgBuffer, op)
			}
		}

	}
}