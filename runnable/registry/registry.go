package registry

import (
	"github.com/shooyaaa/core/network"
)

type Registry struct {
	network.Server
}

func NewRegistry() *Registry {
	r := Registry{
		Server: &network.Tcp{},
	}
	return &r
}
