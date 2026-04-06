package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	tokenDurationStr := getEnv("TOKEN_DURATION", "15m")
	tokenDuration, err := time.ParseDuration(tokenDurationStr)
	if err != nil {
		log.Printf("Invalid TOKEN_DURATION '%s', defaulting to 15m: %v", tokenDurationStr, err)
		tokenDuration = time.Minute * 15
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "ledgeguard.db"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-it"),
		TokenDuration: tokenDuration,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
