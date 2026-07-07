package token_bucket

type RateLimiter interface {
	Allow(userID string) bool
}
