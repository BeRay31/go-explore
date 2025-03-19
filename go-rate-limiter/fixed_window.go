package rate_limiter

import (
	"time"
)

func NewFixedRateLimiter(size int, window time.Duration) *FixedRateLimiter {
	rl := &FixedRateLimiter{
		size:    size,
		window:  window,
		clients: map[string]*RLProps{},
	}
	ticker := time.NewTicker(rl.window)
	go func() {
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *FixedRateLimiter) cleanup() {
	// clean each clients
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, client := range rl.clients {
		go func() {
			client.mu.Lock()
			now := time.Now()
			for i := 0; i < len(client.reqs); i++ {
				if now.Sub(client.reqs[i]) > rl.window {
					client.reqs = client.reqs[i+1:]
				}
			}
			client.mu.Unlock()
		}()
	}
}

func (rl *FixedRateLimiter) Allow(clientId string) bool {
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
	if len(client.reqs) < rl.size {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.reqs = append(client.reqs, reqTime)
		return true
	}
	return false
}
