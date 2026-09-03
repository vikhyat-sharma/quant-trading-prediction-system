package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

// AuthMiddleware validates JWT tokens and stores claims in request context.
// Claims are never stored in mutable request headers.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := util.VerifyJWT(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, contextKeyUserEmail, claims.Email)
		ctx = context.WithValue(ctx, contextKeyUserRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware populates context if a valid token is present.
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				if claims, err := util.VerifyJWT(parts[1]); err == nil {
					ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
					ctx = context.WithValue(ctx, contextKeyUserEmail, claims.Email)
					ctx = context.WithValue(ctx, contextKeyUserRole, claims.Role)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware requires admin role; must be applied after AuthMiddleware.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !util.IsAdminRole(ContextUserRole(r.Context())) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RateLimiter implements a sliding-window in-memory rate limiter.
type RateLimiter struct {
	mu            sync.Mutex
	requestCounts map[string][]time.Time
	maxRequests   int
	window        time.Duration
	stop          chan struct{}
}

// NewRateLimiter creates a rate limiter. Call Stop() when the server shuts down.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requestCounts: make(map[string][]time.Time),
		maxRequests:   maxRequests,
		window:        window,
		stop:          make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

// Stop terminates the background cleanup goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stop)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			rl.mu.Lock()
			for ip, times := range rl.requestCounts {
				filtered := times[:0]
				for _, t := range times {
					if now.Sub(t) < rl.window {
						filtered = append(filtered, t)
					}
				}
				if len(filtered) == 0 {
					delete(rl.requestCounts, ip)
				} else {
					rl.requestCounts[ip] = filtered
				}
			}
			rl.mu.Unlock()
		case <-rl.stop:
			return
		}
	}
}

// Middleware returns an http.Handler middleware for rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		rl.mu.Lock()
		times := rl.requestCounts[ip]
		filtered := times[:0]
		for _, t := range times {
			if now.Sub(t) < rl.window {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) >= rl.maxRequests {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		rl.requestCounts[ip] = append(filtered, now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// generateRequestID returns a cryptographically random 16-byte hex string.
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback: use nanosecond timestamp (not ideal but won't panic)
		return hex.EncodeToString([]byte(time.Now().String()))[:16]
	}
	return hex.EncodeToString(b)
}

// RecoveryMiddleware recovers from panics and returns 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered request_id=%s: %v",
					r.Header.Get("X-Request-ID"), err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
