package rate_limiter

import (
	"time"
)

// sliding window rate limiter
func NewSlidingRateLimiter(size int, window time.Duration) *SlidingRateLimiter {
	return &SlidingRateLimiter{
		size:    size,
		window:  window,
		clients: map[string]*RLProps{},
	}
}

// Sliding window part
func (rl *SlidingRateLimiter) cleanup(clientId string, timeNow time.Time) {
	// clean each clients
	rl.clients[clientId].mu.Lock()
	defer rl.clients[clientId].mu.Unlock()
	for i := 0; i < len(rl.clients[clientId].reqs); i++ {
		if timeNow.Sub(rl.clients[clientId].reqs[i]) > rl.window {
			rl.clients[clientId].reqs = rl.clients[clientId].reqs[i+1:]
		} else {
			break // as time will be in sorted order
		}
	}
}

func (rl *SlidingRateLimiter) Allow(clientId string) bool {
	if rl.clients[clientId] == nil { // client not yet established
		rl.mu.Lock()
		rl.clients[clientId] = &RLProps{
			reqs:     []time.Time{},
			clientId: clientId,
		}
		defer rl.mu.Unlock()
	}

	client := rl.clients[clientId]
	reqTime := time.Now()
	rl.cleanup(client.clientId, reqTime)
	if len(client.reqs) < rl.size {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.reqs = append(client.reqs, reqTime)
		return true
	}
	return false
}
