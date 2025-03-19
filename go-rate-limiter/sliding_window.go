package rate_limiter

import (
	"time"
)

// sliding window rate limiter
func NewSlidingRateLimiter(size int, window time.Duration) *SlidingRateLimiter {
	return &SlidingRateLimiter{
		size:    size,
		window:  window,
		Clients: map[string]*RLProps{},
	}
}

// Sliding window part
func (rl *SlidingRateLimiter) Cleanup(clientId string, timeNow time.Time) {
	// clean each Clients
	rl.Clients[clientId].mu.Lock()
	defer rl.Clients[clientId].mu.Unlock()
	for i := 0; i < len(rl.Clients[clientId].Reqs); i++ {
		if timeNow.Sub(rl.Clients[clientId].Reqs[i]) > rl.window {
			rl.Clients[clientId].Reqs = rl.Clients[clientId].Reqs[i+1:]
		} else {
			break // as time will be in sorted order
		}
	}
}

func (rl *SlidingRateLimiter) Allow(clientId string) bool {
	if rl.Clients[clientId] == nil { // client not yet established
		rl.mu.Lock()
		rl.Clients[clientId] = &RLProps{
			Reqs:     []time.Time{},
			clientId: clientId,
		}
		defer rl.mu.Unlock()
	}

	client := rl.Clients[clientId]
	reqTime := time.Now()
	rl.Cleanup(client.clientId, reqTime)
	if len(client.Reqs) < rl.size {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.Reqs = append(client.Reqs, reqTime)
		return true
	}
	return false
}
