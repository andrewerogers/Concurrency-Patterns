package patterns

import (
	"sync"
)

const(
	MAX_WORKERS = 5
)

// FanOut will start MAX_WORKERS number of read pipelines
// for the designated in channel
func StartWorkers(in <-chan int, job func(in <-chan int) <-chan int) []<- chan int {
	var workers []<-chan int
	workerCount := 1
	for workerCount <= MAX_WORKERS {
		workers = append(workers, job(in))
		workerCount++
	}
	return workers
}

// Merge will read from all worker pipelines in cs and write to one out channel
func Merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Read from worker channel and write to single out channel
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Wait for all worker channels to close, then close output channel
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
