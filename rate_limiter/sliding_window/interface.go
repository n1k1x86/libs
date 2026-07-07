package sliding_window

type RateLimiter interface {
	Allow(userID string) bool
}
