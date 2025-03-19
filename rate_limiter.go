package rate_limiter

import "time"

// sliding window rate limiter
type RateLimiter struct {
	size   int
	window time.Duration
	reqs   []time.Time
}

func NewRateLimiter(size int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		size:   size,
		window: window,
	}
	ticker := time.NewTicker(rl.window)
	go func() {
		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			}
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	now := time.Now()
	for i := 0; i < len(rl.reqs); i++ {
		if now.Sub(rl.reqs[i]) > rl.window {
			rl.reqs = rl.reqs[i+1:]
		}
	}
}

func (rl *RateLimiter) Allow() bool {
	reqTime := time.Now()
	if len(rl.reqs) < rl.size {
		rl.reqs = append(rl.reqs, reqTime)
		return true
	}
	return false
}
