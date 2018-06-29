package main

import (
	"github.com/shooyaaa/connector"
)

func main() {
	server := connector.HttpServer{"./", ":3333"}
	server.Run()
}
