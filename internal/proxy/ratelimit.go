package proxy

import (
	"sync"
	"time"
)



type RateLimiter struct {
	limit int
	window time.Duration
	requests map[string]int
	mutex sync.Mutex
}


func NewRateLimiter(limit int , window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit: limit,
		window: window,
		requests: make(map[string]int),
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()


	count := r.requests[ip]
	if count > 100 {
		return false
	}

	r.requests[ip] = count + 1
	time.AfterFunc(time.Minute, func() {
		r.mutex.Lock()
		r.requests[ip]--
		if r.requests[ip] <= 0 {
			delete(r.requests, ip)
		}

		r.mutex.Unlock()
	})

	return true
}
