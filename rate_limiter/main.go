package rate_limiter

import (
	"sync"
	"time"
)

// sliding window rate limiter
type RateLimiter struct {
	size   int
	window time.Duration
	reqs   []time.Time
	mu     sync.Mutex
}

func NewRateLimiter(size int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		size:   size,
		window: window,
	}
	ticker := time.NewTicker(rl.window)
	go func() {
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	now := time.Now()
	for i := 0; i < len(rl.reqs); i++ {
		if now.Sub(rl.reqs[i]) > rl.window {
			rl.reqs = rl.reqs[i+1:]
		}
	}
	rl.mu.Unlock()
}

func (rl *RateLimiter) Allow() bool {
	reqTime := time.Now()
	if len(rl.reqs) < rl.size {
		rl.mu.Lock()
		rl.reqs = append(rl.reqs, reqTime)
		rl.mu.Unlock()
		return true
	}
	return false
}
