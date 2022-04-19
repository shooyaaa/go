package types

type Game interface {
	Play([]OpWithSession)
}
