package types

import "github.com/shooyaaa/core/op"

type Game interface {
	Play([]op.Op)
}
