package game

import (
	"github.com/shooyaaa/types"
)
type Snake struct {
	r Room	
}

type SnakeData struct {
	Player
}

func (s *Snake) SetRoom(r Room) {
	s.r = r
}

func (s *Snake) HandleOps(ops []Op) {

}

func (s *Snake) Sync() {

}

func (s *Snake) GameData() {
	return SnakeData{X : 0, Y : 0, Blood : 100, Score : 0}
}