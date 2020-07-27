package types

type Game interface {
	HandleOps(r *Room)
	Sync(r *Room) []Op
	GameData() *Player
}
