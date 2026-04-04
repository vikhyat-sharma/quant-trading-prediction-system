package config

import "github.com/spf13/viper"

type Config struct {
	Port        string
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable")
	viper.AutomaticEnv()

	config := &Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	return config, nil
}
