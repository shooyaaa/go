package main

import (
	"fmt"
)

type I interface {
	Flush()
}

type S struct {
}

func (s *S) Flush() {
	fmt.Println("in flush")
}

type F func(int, int)

func main() {
	s := S{}
	s.Flush()
}
