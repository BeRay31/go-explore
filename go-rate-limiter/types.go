package rate_limiter

import (
	"sync"
	"time"
)

type RLProps struct {
	Reqs     []time.Time
	mu       sync.Mutex
	clientId string
}

type FixedRateLimiter struct {
	size    int
	window  time.Duration
	Clients map[string]*RLProps
	mu      sync.Mutex
}

type SlidingRateLimiter struct {
	size    int
	window  time.Duration
	Clients map[string]*RLProps
	mu      sync.Mutex
}

type RateLimiter interface {
	Allow(clientId string) bool
}
