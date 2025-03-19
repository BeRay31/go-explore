package rate_limiter

import (
	"sync"
	"time"
)

type RLProps struct {
	reqs     []time.Time
	mu       sync.Mutex
	clientId string
}

type RateLimiter struct {
	size    int
	window  time.Duration
	clients map[string]*RLProps
	mu      sync.Mutex
}

type FixedRateLimiter RateLimiter

type SlidingRateLimiter RateLimiter
