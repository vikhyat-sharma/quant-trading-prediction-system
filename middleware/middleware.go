package middleware

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
)

const maxRequestBodyBytes = 1 << 20 // 1 MiB

// LoggingMiddleware logs HTTP requests with a unique request ID.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		log.Printf("request_id=%s method=%s path=%s status=%d duration=%s",
			requestID, r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware adds CORS headers. The allowed origin is read from constants.
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

// ContentTypeMiddleware sets JSON content type for API responses.
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

// BodyLimitMiddleware rejects request bodies larger than maxRequestBodyBytes.
func BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		next.ServeHTTP(w, r)
	})
}

// LivenessHandler returns 200 OK if the process is alive.
// It does NOT check external dependencies.
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// ReadinessHandlerFunc returns a handler that checks DB connectivity.
// The server is ready only when the database is reachable.
func ReadinessHandlerFunc(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
		if err := database.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "unavailable",
				"reason": "database unreachable",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	}
}
