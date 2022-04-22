package network

import (
	"github.com/shooyaaa/core/session"
)

type Server interface {
	Listen(addr string) error
	Accept() *session.Session
	Close() error
}
