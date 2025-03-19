package main

import (
	"fmt"
	"math/rand"
	"time"

	"poc-example.com/rate_limiter"
)

func main() {
	rl := rate_limiter.NewRateLimiter(30, 3*time.Second)
	// Test without concurrency
	// for i := 1; i <= 7; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Printf("Request %d: at %v, Allowed? %v\n", i, time.Now().Second(), rl.Allow())
	// }
	// Test with concurrency
	size := 100
	jobs := make(chan int, size)
	done := make(chan bool, size)
	for i := 1; i <= size; i++ { // worker
		go func() {
			for j := range jobs {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				isAllowed := rl.Allow()
				fmt.Printf("Request %d: at %v, Allowed? %v\n", j, time.Now().Second(), rl.Allow())
				done <- isAllowed
			}
		}()
	}
	for i := 1; i <= size; i++ {
		jobs <- i
	}
	close(jobs)
	allowed := 0
	for i := 0; i < size; i++ {
		isAllowed := <-done
		if isAllowed {
			allowed += 1
		}
	}
	fmt.Println("Total allowed requests:", allowed)
}
