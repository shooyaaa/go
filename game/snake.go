package game

import (
	"github.com/shooyaaa/types"
)

type Snake struct {
	r types.Room
}

type SnakeData struct {
	types.Player
}

func (s *Snake) SetRoom(r types.Room) {
	s.r = r
}

func (s *Snake) HandleOps(ops []types.Op) {

}

func (s *Snake) Sync() {

}

func (s *Snake) GameData() SnakeData {
	data := SnakeData{}
	data.Blood = 100
	return data
}
