package token_bucket

import (
	"sync"
	"time"
)

type limiter struct {
	mu       sync.Mutex
	rate     float64
	capacity float64
	buckets  map[string]*bucket
}

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

func New(rate, capacity float64) RateLimiter {
	return &limiter{
		rate:     rate,
		capacity: capacity,
		buckets:  map[string]*bucket{},
	}
}

func (l *limiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	b, ok := l.buckets[userID]
	if !ok {
		b = &bucket{
			tokens:     l.capacity,
			lastRefill: now,
		}
		l.buckets[userID] = b
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * l.rate

	if b.tokens >= l.capacity {
		b.tokens = l.capacity
	}

	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}

	b.tokens -= 1
	return true
}
