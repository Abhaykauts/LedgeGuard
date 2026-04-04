package config

import (
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Port          string
	DBPath        string
	JWTSecret     string
	TokenDuration time.Duration
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "ledgeguard.db"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-it"),
		TokenDuration: time.Minute * 15,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
