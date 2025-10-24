package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-starter/internal/logger"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting per IP address
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      int
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

// getLimiter returns a rate limiter for the given IP address
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// cleanup removes old entries from the limiters map
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// In production, you might want to track last access time
		// For now, we clear all limiters periodically
		rl.limiters = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a middleware that rate limits requests by IP
func RateLimitMiddleware(rps, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rps, burst)

	// Start cleanup goroutine
	go limiter.cleanup()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := getClientIP(r)

			// Get or create limiter for this IP
			ipLimiter := limiter.getLimiter(ip)

			// Check if request is allowed
			if !ipLimiter.Allow() {
				// Log rate limit exceeded
				logger.FromContext(r.Context()).Warn("rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)

				// Calculate retry-after duration
				reservation := ipLimiter.Reserve()
				if !reservation.OK() {
					reservation.Cancel()
					w.Header().Set("Retry-After", "60")
				} else {
					delay := reservation.Delay()
					reservation.Cancel()
					w.Header().Set("Retry-After", delay.String())
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"too many requests","message":"rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for reverse proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
