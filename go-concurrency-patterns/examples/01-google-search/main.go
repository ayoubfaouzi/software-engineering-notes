package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string
type Search func(query string) Result

var (
	Web1   = fakeSearch("web1")
	Web2   = fakeSearch("web2")
	Image1 = fakeSearch("image1")
	Image2 = fakeSearch("image2")
	Video1 = fakeSearch("video1")
	Video2 = fakeSearch("video2")
)

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep((time.Duration(rand.Intn(100)) * time.Millisecond))
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// How do we avoid discarding result from the slow server?
// We duplicates too many instances, and perfor parallel requests.
func First(query string, replicas ...Search) Result {

	c := make(chan Result)

	for i := range replicas {
		go func(idx int) {
			c <- replicas[idx](query)
		}(i)
	}

	// First function always wait for 1 time after receiving the result.
	return <-c

}

// Don't wait for the slowest server.
func Google(query string) []Result {

	c := make(chan Result)
	var results []Result

	// each search is performed in a goroutine.
	go func() {
		c <- First(query, Web1, Web2)
	}()
	go func() {
		c <- First(query, Image1, Image2)
	}()
	go func() {
		c <- First(query, Video1, Video2)
	}()

	timeout := time.After(100 * time.Millisecond)

	for range 3 {
		select {
		case r := <-c:
			results = append(results, r)
		case <-timeout:
			fmt.Println("timeout")
			return results
		}
	}
	return results

}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
