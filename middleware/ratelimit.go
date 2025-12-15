package middleware

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimitConfig holds rate limit settings
type RateLimitConfig struct {
	RPS   rate.Limit
	Burst int
}

// NewRateLimiter creates a shared rate limiter
func NewRateLimiter(cfg RateLimitConfig) *rate.Limiter {
	return rate.NewLimiter(cfg.RPS, cfg.Burst)
}

// RateLimitMiddleware applies rate limiting to an HTTP handler
func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("rate limit exceeded"))
				fmt.Println("Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}