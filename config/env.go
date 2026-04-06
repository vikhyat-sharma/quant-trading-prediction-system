package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
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
	viper.SetDefault(constants.EnvKeyPort, constants.DefaultPort)
	viper.SetDefault(constants.EnvKeyDatabaseURL, constants.DefaultDatabaseURL)
	viper.SetDefault(constants.EnvKeyEnvironment, constants.DefaultEnvironment)
	viper.SetDefault(constants.EnvKeyLogLevel, constants.DefaultLogLevel)
	viper.AutomaticEnv()

	config := &Config{
		Port:        viper.GetString(constants.EnvKeyPort),
		DatabaseURL: viper.GetString(constants.EnvKeyDatabaseURL),
		Environment: viper.GetString(constants.EnvKeyEnvironment),
		LogLevel:    viper.GetString(constants.EnvKeyLogLevel),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
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
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf(constants.ErrMsgLogLevelInvalid)
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
