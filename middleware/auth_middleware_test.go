package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

func TestAuthMiddlewareSetsNumericUserIDHeader(t *testing.T) {
	token, err := util.GenerateJWT(7, "user@example.com", "user")
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-User-ID"); got != "7" {
			t.Fatalf("expected X-User-ID to be numeric user id, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestGetJWTSecretUsesEnvironmentVariable(t *testing.T) {
	t.Setenv("JWT_SECRET", "env-secret")

	if secret := util.GetJWTSecret(); secret != "env-secret" {
		t.Fatalf("expected JWT secret from environment, got %q", secret)
	}
}
