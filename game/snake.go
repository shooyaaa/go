package game

import (
	"log"

	"github.com/shooyaaa/types"
)

const (
	Move = 1
	Sync = 2
)

type Snake struct {
}

type SnakeData struct {
	types.Player
}

func (s Snake) HandleOps(r *types.Room) {
	for _, msg := range r.GetMsgBuffer() {
		session := msg.GetId()
		gameData, err := r.GetMember(session)
		if err != nil {
			log.Printf("Error while HandleOps %v", err)
			continue
		}
		switch msg.Type {
		case Move:
			x := int(msg.Data["X"])
			gameData.X = x
			y := int(msg.Data["Y"])
			gameData.Y = y
		}
	}
}

func (s Snake) Sync(r *types.Room) []types.Op {
	ops := make([]types.Op, r.MemberCount())
	count := 0
	for session, gameData := range r.AllMembers() {
		dict := make(map[string]float64)
		dict["x"] = float64(gameData.X)
		dict["y"] = float64(gameData.Y)
		dict["id"] = float64(session.Id)
		ops[count] = types.Op{
			Type: Sync,
			Ts:   r.FrameTime,
			Data: dict,
		}
		count++
	}
	return ops
}

func (s Snake) GameData() *types.Player {
	data := types.Player{}
	data.Blood = 100
	return &data
}
