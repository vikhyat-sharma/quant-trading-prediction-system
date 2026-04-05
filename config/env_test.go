package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")

	config, err := LoadConfig()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if config == nil {
		t.Errorf("expected config, got nil")
	}

	// Defaults should be set
	if config.Port == "" {
		t.Errorf("expected Port to have a default value")
	}

	if config.DatabaseURL == "" {
		t.Errorf("expected DatabaseURL to have a default value")
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")

	config, err := LoadConfig()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if config.Port != "9000" {
		t.Errorf("expected Port 9000, got %s", config.Port)
	}

	if config.DatabaseURL != "postgres://test:test@localhost/testdb" {
		t.Errorf("expected specific DatabaseURL, got %s", config.DatabaseURL)
	}

	// Cleanup
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
}

func TestLoadConfig_StructFields(t *testing.T) {
	config, err := LoadConfig()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that Config struct has the expected fields
	if config.Port == "" && config.DatabaseURL == "" {
		t.Errorf("expected Config to have Port and DatabaseURL fields populated")
	}
}
