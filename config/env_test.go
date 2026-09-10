package config

import (
	"testing"
)

func TestLoadConfig_MissingDatabaseURL_ReturnsError(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "info")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is empty, got nil")
	}
}

func TestLoadConfig_WithValidEnvVars(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != "9000" {
		t.Errorf("expected Port 9000, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://test:test@localhost/testdb" {
		t.Errorf("unexpected DatabaseURL: %s", cfg.DatabaseURL)
	}
}

func TestLoadConfig_InvalidPort_ReturnsError(t *testing.T) {
	t.Setenv("PORT", "notaport")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "info")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid PORT, got nil")
	}
}

func TestLoadConfig_InvalidEnvironment_ReturnsError(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "unknown")
	t.Setenv("LOG_LEVEL", "info")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid ENVIRONMENT, got nil")
	}
}

func TestLoadConfig_ProductionRequiresJWTSecret(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("JWT_SECRET", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is missing in production, got nil")
	}
}

func TestLoadConfig_ProductionWithJWTSecret_Succeeds(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("JWT_SECRET", "a-strong-secret-that-is-long-enough-32c")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !cfg.IsProduction() {
		t.Error("expected IsProduction() to return true")
	}
}

func TestLoadConfig_DevelopmentWithoutJWTSecret_Succeeds(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("JWT_SECRET", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error in development without JWT_SECRET, got %v", err)
	}
	if !cfg.IsDevelopment() {
		t.Error("expected IsDevelopment() to return true")
	}
}
