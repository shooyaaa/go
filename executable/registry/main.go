package main

import (
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/runnable/registry"
)

func main() {
	r := registry.NewRegistry()
	r.Listen(config.RegistryRedisAddress)
	r.Accept()

	defer r.Close()
}
