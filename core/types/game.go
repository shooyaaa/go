package types

import "github.com/shooyaaa/core/session"

type Game interface {
	Play([]session.Op)
}
