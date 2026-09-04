package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

func makeToken(t *testing.T, userID int, email, role string) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-32chars")
	token, err := util.GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	return token
}

func TestAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	token := makeToken(t, 7, "user@example.com", "user")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := ContextUserID(r.Context()); got != 7 {
			t.Errorf("expected userID 7 in context, got %d", got)
		}
		if got := ContextUserEmail(r.Context()); got != "user@example.com" {
			t.Errorf("expected email in context, got %q", got)
		}
		if got := ContextUserRole(r.Context()); got != "user" {
			t.Errorf("expected role in context, got %q", got)
		}
		// Claims must NOT be in headers (security: headers are client-controlled)
		if r.Header.Get("X-User-ID") != "" {
			t.Error("X-User-ID header should not be set by AuthMiddleware")
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingHeader_Returns401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken_Returns401(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-32chars")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()
	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_MalformedHeader_Returns401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer token")
	w := httptest.NewRecorder()
	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	})).ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetJWTSecret_UsesEnvironmentVariable(t *testing.T) {
	t.Setenv("JWT_SECRET", "env-secret")
	if secret := util.GetJWTSecret(); secret != "env-secret" {
		t.Fatalf("expected JWT secret from environment, got %q", secret)
	}
}

func TestGetJWTSecret_EmptyWhenNotSet(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if secret := util.GetJWTSecret(); secret != "" {
		t.Fatalf("expected empty secret when env var not set, got %q", secret)
	}
}

func TestRateLimiter_Stop_DoesNotPanic(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	// Stop must be safe to call and must terminate the goroutine
	done := make(chan struct{})
	go func() {
		rl.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("RateLimiter.Stop() did not return within 2s")
	}
}

func TestRateLimiter_ExceedsLimit_Returns429(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	defer rl.Stop()

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after limit exceeded, got %d", w.Code)
	}
}

func TestRecoveryMiddleware_PanicReturns500(t *testing.T) {
	handler := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 after panic, got %d", w.Code)
	}
}

func TestAdminMiddleware_NonAdmin_Returns403(t *testing.T) {
	token := makeToken(t, 1, "user@example.com", "user")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Chain: AuthMiddleware → AdminMiddleware → handler
	chain := AuthMiddleware(AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for non-admin")
	})))
	chain.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAdminMiddleware_Admin_Passes(t *testing.T) {
	token := makeToken(t, 1, "admin@example.com", "admin")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	chain := AuthMiddleware(AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	chain.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin, got %d", w.Code)
	}
}
