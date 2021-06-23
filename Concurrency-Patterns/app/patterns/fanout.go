package patterns

import (
	"fmt"
	"runtime"
	"time"
)

// FAN-OUT Concurrency pattern
// Works kind of like a load balancer to manage server load
//		# One single read channel
//		# Multiple processing channels to consume read
//		# Process channels write to a single output channel
func RunFanOut() {
	start := time.Now()
	input := read(1,2,3,4,5,6,7,8,9,10)

	workerPipelines := StartWorkers(input, sq)
	// read from started worker pipelines in one single output channel
	for square := range Merge(workerPipelines) {
		fmt.Println(square)
	}

	elapsed := time.Since(start)
	fmt.Println("elapsed: ", elapsed)
	fmt.Println("routines: ", runtime.NumGoroutine())
}

// read returns an integer channel and writes
// asynchronously to the channel for all numbers in nums
func read(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// sq reads integers from an in channel and writes the square
// asynchronously to a returned out channel
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n*n
			time.Sleep(150 * time.Millisecond) // simulate work time
		}
		close(out)
	}()
	return out
}
