package rate_limiter_test

import (
	"sync"
	"testing"
	"time"

	"beray-explore.com/rate_limiter"
)

func TestAllowSequential(t *testing.T) {
	rl := rate_limiter.NewRateLimiter(2, 2*time.Second)
	clientId := "client-1"

	if !rl.Allow(clientId) {
		t.Error("Expected 1st request to be allowed")
	}
	if !rl.Allow(clientId) {
		t.Error("Expected 2nd request to be allowed")
	}
	if rl.Allow(clientId) {
		t.Error("Expected 3rd request to be denied")
	}

	time.Sleep(3 * time.Second)

	if !rl.Allow(clientId) {
		t.Error("Expected 4th request to be allowed after window reset")
	}
}

func TestAllowConcurrent(t *testing.T) {
	rl := rate_limiter.NewRateLimiter(2, 3*time.Second)
	clientId := "client-1"
	var wg sync.WaitGroup
	results := make(chan bool, 4)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- rl.Allow(clientId)
		}()
	}

	wg.Wait()
	close(results)

	allowedCount := 0
	for result := range results {
		if result {
			allowedCount++
		}
	}

	if allowedCount != 2 {
		t.Errorf("Expected 2 allowed requests, got %d", allowedCount)
	}
}
