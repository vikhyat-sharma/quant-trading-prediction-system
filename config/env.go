package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
)

// Config holds all configuration for the application.
type Config struct {
	Port        string
	DatabaseURL string
	Environment string
	LogLevel    string
	JWTSecret   string
	CORSOrigin  string

	// Server timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LoadConfig loads configuration from environment variables with validation.
func LoadConfig() (*Config, error) {
	viper.SetDefault(constants.EnvKeyPort, constants.DefaultPort)
	viper.SetDefault(constants.EnvKeyEnvironment, constants.DefaultEnvironment)
	viper.SetDefault(constants.EnvKeyLogLevel, constants.DefaultLogLevel)
	viper.SetDefault(constants.EnvKeyCORSOrigin, constants.CORSAllowOrigin)
	viper.AutomaticEnv()

	cfg := &Config{
		Port:         viper.GetString(constants.EnvKeyPort),
		DatabaseURL:  viper.GetString(constants.EnvKeyDatabaseURL),
		Environment:  viper.GetString(constants.EnvKeyEnvironment),
		LogLevel:     viper.GetString(constants.EnvKeyLogLevel),
		JWTSecret:    viper.GetString(constants.EnvKeyJWTSecret),
		CORSOrigin:   viper.GetString(constants.EnvKeyCORSOrigin),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration is present and valid.
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf(constants.ErrMsgPortCannotBeEmpty)
	}
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf(constants.ErrMsgPortMustBeValidNumber+": %w", err)
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf(constants.ErrMsgDatabaseURLCannotBeEmpty)
	}

	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[c.Environment] {
		return fmt.Errorf(constants.ErrMsgEnvironmentInvalid)
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf(constants.ErrMsgLogLevelInvalid)
	}

	// JWT_SECRET is required in staging and production.
	if c.Environment != "development" && c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required in %s environment", c.Environment)
	}

	return nil
}

// IsProduction returns true if the environment is production.
func (c *Config) IsProduction() bool { return c.Environment == "production" }

// IsDevelopment returns true if the environment is development.
func (c *Config) IsDevelopment() bool { return c.Environment == "development" }
