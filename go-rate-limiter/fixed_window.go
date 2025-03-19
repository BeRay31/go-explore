package rate_limiter

import (
	"time"
)

func NewFixedRateLimiter(size int, window time.Duration) *FixedRateLimiter {
	rl := &FixedRateLimiter{
		size:    size,
		window:  window,
		Clients: map[string]*RLProps{},
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
	// clean each Clients
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, client := range rl.Clients {
		go func() {
			client.mu.Lock()
			now := time.Now()
			for i := 0; i < len(client.Reqs); i++ {
				if now.Sub(client.Reqs[i]) > rl.window {
					client.Reqs = client.Reqs[i+1:]
				}
			}
			client.mu.Unlock()
		}()
	}
}

func (rl *FixedRateLimiter) Allow(clientId string) bool {
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
	if len(client.Reqs) < rl.size {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.Reqs = append(client.Reqs, reqTime)
		return true
	}
	return false
}
