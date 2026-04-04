package http

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ipLimiter holds a rate limiter per IP address.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimitStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
}

var loginRateLimiter = &rateLimitStore{
	limiters: make(map[string]*ipLimiter),
}

func init() {
	// Periodically clean up stale entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			loginRateLimiter.cleanup()
		}
	}()
}

func (s *rateLimitStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.limiters[ip]
	if !ok {
		// Allow 5 login attempts per minute per IP, burst of 5
		entry = &ipLimiter{limiter: rate.NewLimiter(rate.Every(time.Minute/5), 5)}
		s.limiters[ip] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

func (s *rateLimitStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for ip, entry := range s.limiters {
		if entry.lastSeen.Before(cutoff) {
			delete(s.limiters, ip)
		}
	}
}

// realIP extracts the real client IP, respecting X-Forwarded-For from trusted proxies.
func realIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		// Take the first (leftmost) IP which is the original client
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	if fwd := r.Header.Get("X-Real-Ip"); fwd != "" {
		return strings.TrimSpace(fwd)
	}
	// Strip port from RemoteAddr
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// withLoginRateLimit wraps a handleFunc with per-IP rate limiting.
// Returns 429 Too Many Requests when the limit is exceeded.
func withLoginRateLimit(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, data *requestContext) (int, error) {
		ip := realIP(r)
		limiter := loginRateLimiter.get(ip)
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "60")
			return http.StatusTooManyRequests, errTooManyRequests
		}
		return fn(w, r, data)
	}
}

var errTooManyRequests = &rateLimitError{}

type rateLimitError struct{}

func (e *rateLimitError) Error() string {
	return "too many login attempts, please try again later"
}
