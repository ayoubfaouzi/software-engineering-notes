package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Msg struct {
	str  string
	wait chan bool
}

func boring(msg string) <-chan string { // Return receive-only channel of strings

	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c

}

func boringWithOrder(msg string) <-chan Msg {
	c := make(chan Msg)
	waitForIt := make(chan bool) // shared between all messages
	go func() {
		for i := 0; ; i++ {
			c <- Msg{
				str:  fmt.Sprintf("%s %d", msg, i),
				wait: waitForIt,
			}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

			// The code waits until the value to be received.
			<-waitForIt
		}

	}()
	return c
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)

	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}

	}()
	return c
}

func fanInWithOrder(inputs ...<-chan Msg) <-chan Msg {

	c := make(chan Msg)

	for i := range inputs {
		input := inputs[i]
		go func() {
			for {
				c <- <-input
			}
		}()
	}
	return c
}

// Rewrite our original fanIn function. Only one goroutine is needed.
func fanInWithSelect(input1, input2 <-chan string) <-chan string {

	c := make(chan string)

	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s

			}
		}

	}()

	return c
}

func main() {
	// No order
	c := fanIn(boring("Joe"), boring("Ann"))
	for range 10 {
		fmt.Println(<-c)
	}
	fmt.Println("You're both boring. I'm leaving.")

	// Force order.
	c2 := fanInWithOrder(boringWithOrder("Joe"), boringWithOrder("Ann"))
	for range 5 {
		msg1 := <-c2 // block until the message is read
		fmt.Println(msg1.str)
		msg2 := <-c2
		fmt.Println(msg2.str)

		// each go routine have to wait
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're both boring. I'm leaving.")

	// Using select and timeouts.
	c3 := fanInWithSelect(boring("Joe"), boring("Ann"))
	timeout := time.After(3 * time.Second)
	for {
		select {
		case s := <-c3:
			fmt.Println(s)
		case <-timeout:
			fmt.Println("global timeout reached.")
			return
		case <-time.After(1 * time.Second): // returns a channel after specified interval, this is a timeout for each iteration in the loop.
			fmt.Println("You're too slow.")
			return
		}
	}
}
