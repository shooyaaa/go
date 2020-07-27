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
	rooms map[types.ID]*types.Room
	uuid  types.UUID
}

func (rm *roomManager) Init() {
	rm.rooms = make(map[types.ID]*types.Room)
	rm.uuid = &types.Simple{}
}

func (rm *roomManager) Join(id types.ID, playerId *types.Session) error {
	r, ok := rm.rooms[id]
	if !ok {
		return errors.New("Room %d not exists")
	}
	r.Add(playerId)
	return nil
}

func (rm *roomManager) Add() (types.ID, *types.Room) {
	id := rm.uuid.NewUUID()
	room := types.Room{}
	room.Init()
	rm.rooms[id] = &room
	return id, rm.rooms[id]
}

func (rm *roomManager) Get(id types.ID) (*types.Room, error) {
	r, ok := rm.rooms[id]
	if !ok {
		return nil, errors.New("Room %d not exists")
	}
	return r, nil
}
