package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		results <- job * 2
	}
}

func main() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	var wg sync.WaitGroup

	// Fan-out: 3 workers
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for i := 1; i <= 5; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)

	// Fan-in: collect results
	for res := range results {
		fmt.Println("Result:", res)
	}
}
