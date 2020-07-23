package types


type Game interface {
	SetRoom(r Room)
	HandleOps(ops []Op)
	Sync()
	GameData() Player
}