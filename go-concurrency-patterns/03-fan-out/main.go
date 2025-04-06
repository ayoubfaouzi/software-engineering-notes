package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, val := range nums {
			c <- val
		}

		close(c)
	}()

	return c
}

func main() {

	data1 := []int{1, 2, 3, 4, 5}
	data2 := []int{10, 20, 30, 40, 50}

	var wg sync.WaitGroup
	ch1 := generator(data1...)
	ch2 := generator(data2...)
	wg.Add(2)

	go func() {
		for val := range ch1 {
			fmt.Printf("Channel1 data: %v\n", val)
		}

		wg.Done()
	}()
	go func() {
		for val := range ch2 {
			fmt.Printf("Channel2 data: %v\n", val)
		}

		wg.Done()
	}()

	wg.Wait()

}
