package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

// AuthMiddleware validates JWT tokens
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

		token := parts[1]
		claims, err := util.VerifyJWT(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store claims in context for later use
		r.Header.Set("X-User-ID", string(rune(claims.UserID)))
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuthMiddleware allows both authenticated and unauthenticated requests
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				claims, err := util.VerifyJWT(token)
				if err == nil {
					r.Header.Set("X-User-ID", string(rune(claims.UserID)))
					r.Header.Set("X-User-Email", claims.Email)
					r.Header.Set("X-User-Role", claims.Role)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Header.Get("X-User-Role")
		if !util.IsAdminRole(userRole) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware implements basic rate limiting
type RateLimiter struct {
	requestCounts map[string][]time.Time
	maxRequests   int
	window        time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requestCounts: make(map[string][]time.Time),
		maxRequests:   maxRequests,
		window:        window,
	}

	// Cleanup old entries periodically
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			for ip := range rl.requestCounts {
				var filtered []time.Time
				for _, t := range rl.requestCounts[ip] {
					if now.Sub(t) < window {
						filtered = append(filtered, t)
					}
				}
				if len(filtered) == 0 {
					delete(rl.requestCounts, ip)
				} else {
					rl.requestCounts[ip] = filtered
				}
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		// Clean old requests
		var filtered []time.Time
		for _, t := range rl.requestCounts[ip] {
			if now.Sub(t) < rl.window {
				filtered = append(filtered, t)
			}
		}

		if len(filtered) >= rl.maxRequests {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.requestCounts[ip] = append(filtered, now)
		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware adds request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)

		log.Printf("[%s] %s %s", requestID, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().Nanosecond()))
}
