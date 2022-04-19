package game

import (
	"github.com/shooyaaa/types"
)

const (
	Move = 1
	Sync = 2
)

type Snake struct {
	players *types.Player
}

type SnakeData struct {
	types.Player
}

func (s Snake) Play(ops []types.OpWithSession) {
	for _, op := range (ops) {

	}
}

func (s Snake) GameData() *types.Player {
	data := types.Player{}
	data.Blood = 100
	return &data
}
