package network

import "github.com/shooyaaa/types"

type Server interface {
	Listen(addr string) error
	Accept() *types.Session
	Close() error
}
