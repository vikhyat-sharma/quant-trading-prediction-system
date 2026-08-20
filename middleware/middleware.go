package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		r = r.WithContext(r.Context())
		r.Header.Set("X-Request-ID", requestID)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		wrapped.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("request_id=%s method=%s path=%s status=%d duration=%s", requestID, r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderAccessControlAllowOrigin, constants.CORSAllowOrigin)
		w.Header().Set(constants.HeaderAccessControlAllowMethods, constants.CORSAllowMethods)
		w.Header().Set(constants.HeaderAccessControlAllowHeaders, constants.CORSAllowHeaders)

		if r.Method == constants.MethodOPTIONS {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ContentTypeMiddleware ensures JSON content type for API responses
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

// HealthHandler serves readiness and liveness endpoints.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func generateRequestID() string {
	return time.Now().UTC().Format("20060102150405") + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
