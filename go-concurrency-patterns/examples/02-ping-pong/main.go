package main

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table // Player grabs the ball.
		ball.hits++
		fmt.Printf("%s %d\n", name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball // pass the ball.
	}
}

func main() {
	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)

	table <- new(Ball) // game on: toss the ball.
	time.Sleep(1 * time.Second)
	<-table // game over, snatch the ball.
	panic("show me the stack")
}
