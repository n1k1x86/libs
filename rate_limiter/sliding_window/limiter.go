package sliding_window

import (
	"container/list"
	"sync"
	"time"
)

type limiter struct {
	window  time.Duration
	mu      sync.Mutex
	records map[string]*list.List
	limit   int
}

func New(limit int, window time.Duration) RateLimiter {
	return &limiter{
		limit:   limit,
		window:  window,
		records: make(map[string]*list.List),
	}
}

func (l *limiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	userList, ok := l.records[userID]
	if !ok {
		userList := list.New()
		l.records[userID] = userList
	}

	for userList.Len() > 0 {
		front := userList.Front()
		ts := front.Value.(time.Time)

		if now.Sub(ts) <= l.window {
			break
		}

		userList.Remove(front)
	}

	if userList.Len() >= l.limit {
		return false
	}

	userList.PushBack(now)

	return true
}
