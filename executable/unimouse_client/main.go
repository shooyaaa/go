package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/shooyaaa/runnable/unimouse"
)

func main() {
	robotgo.MouseSleep = 100

	robotgo.ScrollMouse(10, "up")
	robotgo.ScrollMouse(20, "right")

	unimouse.Connect()
}
