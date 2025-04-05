package main

import "fmt"

const (
	N = 1000
)

func f(left, right chan int) {
	left <- 1 + <-right
}

func main() {
	leftmost := make(chan int)
	left := leftmost
	right := leftmost

	for i := 0; i < N; i++ {
		right = make(chan int) // create a new channel
		go f(left, right)      // create a new goroutine
		left = right
	}

	go func(c chan int) {
		c <- 1
	}(right)
	fmt.Println(<-leftmost)
}
