package rate_limiter_test

import (
	"sync"
	"testing"
	"time"

	rate_limiter "beray-explore.com"
)

func TestCleanup(t *testing.T) {
	rl := rate_limiter.NewSlidingRateLimiter(5, 2*time.Second)
	clientId := "client-1"

	// Add requests to the client
	rl.Allow(clientId)
	rl.Allow(clientId)
	rl.Allow(clientId)

	// Verify that there are 3 requests
	client := rl.Clients[clientId]
	if client != nil && len(client.Reqs) != 3 {
		t.Errorf("Expected 3 requests, got %d", len(rl.Clients[clientId].Reqs))
	}

	// Call cleanup after 3 seconds
	time.Sleep(3 * time.Second)
	rl.Cleanup(clientId, time.Now())
	rl.Allow(clientId)

	// Verify that old requests are removed
	if len(rl.Clients[clientId].Reqs) != 1 {
		t.Errorf("Expected 1 request, got %d", len(rl.Clients[clientId].Reqs))
	}
}

func TestAllow(t *testing.T) {
	rl := rate_limiter.NewSlidingRateLimiter(3, 2*time.Second)
	clientId := "client-2"

	// Allow 3 requests
	if !rl.Allow(clientId) {
		t.Errorf("Expected request to be allowed")
	}
	if !rl.Allow(clientId) {
		t.Errorf("Expected request to be allowed")
	}
	if !rl.Allow(clientId) {
		t.Errorf("Expected request to be allowed")
	}

	// 4th request should be denied
	if rl.Allow(clientId) {
		t.Errorf("Expected request to be denied")
	}

	// Wait for the window to expire
	time.Sleep(3 * time.Second)

	// Now the request should be allowed again
	if !rl.Allow(clientId) {
		t.Errorf("Expected request to be allowed")
	}
}

func TestConcurrency(t *testing.T) {
	rl := rate_limiter.NewSlidingRateLimiter(10, 2*time.Second)
	clientId := "client-3"
	var wg sync.WaitGroup

	// Simulate concurrent requests
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Allow(clientId)
		}()
	}

	wg.Wait()

	// Verify that no more than 10 requests are allowed
	client := rl.Clients[clientId]
	if len(client.Reqs) > 10 {
		t.Errorf("Expected at most 10 requests, got %d", len(client.Reqs))
	}
}
