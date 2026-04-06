package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Port        string
	DatabaseURL string
	Environment string
	LogLevel    string
}

// LoadConfig loads configuration from environment variables with validation
func LoadConfig() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost/quant_trading?sslmode=disable")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.AutomaticEnv()

	config := &Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		Environment: viper.GetString("ENVIRONMENT"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be a valid number: %w", err)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL cannot be empty")
	}

	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}

	if !validEnvs[c.Environment] {
		return fmt.Errorf("ENVIRONMENT must be one of: development, staging, production")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	return nil
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
