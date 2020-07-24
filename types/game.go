package types

type Game interface {
	SetRoom(r Room)
	HandleOps(ops []Op)
	Sync() []Op
	GameData() Player
}
