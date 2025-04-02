package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

func main() {
	joe := boring("Joe")
	ann := boring("Ann")
	for i := 0; i < 5; i++ {
		// Because of the sync nature of channels, the first chan (joe)
		// will block (ann's channel) from executing even though `ann`
		// might be ready to send a value!
		// We can get around that by using the `fan` or `multiplexing` pattern.
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
}
