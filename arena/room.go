package arena

import (
	"github.com/shooyaaa/types"
	"github.com/shooyaaa/uuid"
	"log"
)

type player struct {
	x int
	y int
}

type room struct {
	members map[uuid.ID]player
}

type RoomManager struct {
	rooms   map[int]room
	players map[uuid.ID]int
}

func (r *room) Add(id uuid.ID) {
	r.members[id] = player{
		x: 0,
		y: 0,
	}
}

func (r *room) Update(id uuid.ID, x int, y int) error {
	player, ok := r.members[id]
	if !ok {
		return nil
	}
	player.x = x
	player.y = y
	return nil
}

func (r *room) Handle(msg interface{}) {

}

func (rm *RoomManager) Join(id int, playerId uuid.ID) {
	r, ok := rm.rooms[id]
	if !ok {
		r = room{}
		rm.rooms[id] = r
	}
	r.Add(playerId)
	rm.players[playerId] = id
}

func (rm *RoomManager) add(id int) {
	rm.rooms[id] = room{}
}

func (rm RoomManager) Handle(op types.OpCode, data []byte) error {
	log.Printf("Op code v%", op)

	return nil
}
