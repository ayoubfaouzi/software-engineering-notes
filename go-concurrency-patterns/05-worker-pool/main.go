package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func workerEfficient(id int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup

	for j := range jobs {
		wg.Add(1)

		go func(job int) {
			fmt.Println("worker", id, "started job", j)
			time.Sleep(time.Second)
			fmt.Println("worker", id, "finished job", j)
			results <- j * 2
			wg.Done()
		}(j)
	}
}

func main() {
	const numJobs = 8
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// In this example, we define a fixed 3 workers
	// they receive the `jobs` from the channel jobs
	// we also naming the worker name with `w` variable.
	for w := 1; w <= 3; w++ {
		go workerEfficient(w, jobs, results)
	}

	// Push the jobs.
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	fmt.Println("closed jobs")

	for a := 1; a <= numJobs; a++ {
		<-results
	}
	close(results)

}
