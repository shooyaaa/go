package manager

import (
	"github.com/shooyaaa/types"
)


type RoomManager struct {
	rooms   map[int]types.Room
	players map[types.ID]int
}


func (rm *RoomManager) Join(id int, playerId types.ID) {
	r, ok := rm.rooms[id]
	if !ok {
		r = types.Room{}
		rm.rooms[id] = r
	}
	r.Add(playerId)
	rm.players[playerId] = id
}

func (rm *RoomManager) add(id int) {
	rm.rooms[id] = types.Room{}
}

