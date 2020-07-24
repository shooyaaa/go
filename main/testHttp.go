package main

import "fmt"

type test struct {
	x int
}

func change(t chan *test) {
	temp := <-t
	temp.x = 3333
}

func main() {
	c := make(chan *test)
	tt := test{x: 1}
	go change(c)
	c <- &tt
	fmt.Printf("value %v", tt)
}
