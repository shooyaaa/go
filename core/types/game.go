package types

import "github.com/shooyaaa/core/codec"

type Game interface {
	Play([]codec.Op)
}
