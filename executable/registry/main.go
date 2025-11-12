package main

import (
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/runnable/registry"
)

func main() {
	r := registry.NewRegistry()
	if len(config.RegistryRedisAddress) > 0 {
		r.Listen(config.RegistryRedisAddress[0])
	}
	r.Accept()

	defer r.Close()
}
