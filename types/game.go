package types

type Game interface {
	Sync(r *Room) []Op
	GameData() *Player
}
