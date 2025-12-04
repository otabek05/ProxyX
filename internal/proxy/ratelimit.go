package proxy

import (
	"sync"
	"time"
)



type RateLimiter struct {
	limit int
	window time.Duration
	requests map[string]int
	resetAt map[string]time.Time
	mutex sync.Mutex
}


func NewRateLimiter(limit int , window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit: limit,
		window: window,
		requests: make(map[string]int),
		resetAt:make( map[string]time.Time),
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	if resetTime, exists := r.resetAt[ip]; !exists || now.After(resetTime) {
		r.requests[ip] = 0
		r.resetAt[ip] = now.Add(r.window)
	}

	if r.requests[ip] >= r.limit {
		return  false 
	}


	r.requests[ip]++
	return true
}
