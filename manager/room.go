package manager

import (
	"errors"
	"github.com/shooyaaa/core/session"
	types2 "github.com/shooyaaa/core/types"
	"sync"
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
	rooms map[int64]*types2.Room
	uuid  types2.UUID
}

func (rm *roomManager) Init() {
	rm.rooms = make(map[int64]*types2.Room)
	rm.uuid = &types2.Simple{}
}

func (rm *roomManager) Join(id int64, playerId *session.Session) error {
	r, ok := rm.rooms[id]
	if !ok {
		return errors.New("Room %d not exists")
	}
	r.Add(playerId)
	return nil
}

func (rm *roomManager) Add() (int64, *types2.Room) {
	id := rm.uuid.NewUUID()
	room := types2.Room{}
	room.Init()
	rm.rooms[id] = &room
	return id, rm.rooms[id]
}

func (rm *roomManager) Get(id int64) (*types2.Room, error) {
	r, ok := rm.rooms[id]
	if !ok {
		return nil, errors.New("Room %d not exists")
	}
	return r, nil
}
