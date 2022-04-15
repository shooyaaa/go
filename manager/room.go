package manager

import (
	"errors"
	"sync"

	"github.com/shooyaaa/types"
)

var once sync.Once

var instance roomManager

func RoomManager() *roomManager {
	once.Do(func() {
		instance = roomManager{}
		instance.Init()
	})
	return &instance
}

type roomManager struct {
	rooms map[int64]*types.Room
	uuid  types.UUID
}

func (rm *roomManager) Init() {
	rm.rooms = make(map[int64]*types.Room)
	rm.uuid = &types.Simple{}
}

func (rm *roomManager) Join(id int64, playerId *types.Session) error {
	r, ok := rm.rooms[id]
	if !ok {
		return errors.New("Room %d not exists")
	}
	r.Add(playerId)
	return nil
}

func (rm *roomManager) Add() (int64, *types.Room) {
	id := rm.uuid.NewUUID()
	room := types.Room{}
	room.Init()
	rm.rooms[id] = &room
	return id, rm.rooms[id]
}

func (rm *roomManager) Get(id int64) (*types.Room, error) {
	r, ok := rm.rooms[id]
	if !ok {
		return nil, errors.New("Room %d not exists")
	}
	return r, nil
}
