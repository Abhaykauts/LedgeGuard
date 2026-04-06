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
		log.Println("Note: No .env file found; using system environment variables.")
	}

	tokenDurationStr := getEnv("TOKEN_DURATION", "15m")
	tokenDuration, err := time.ParseDuration(tokenDurationStr)
	if err != nil {
		log.Printf("Warning: Invalid TOKEN_DURATION '%s', defaulting to 15m", tokenDurationStr)
		tokenDuration = time.Minute * 15
	}

	dbPath := getEnv("DB_PATH", "ledgeguard.db")
	port := getEnv("PORT", "8080")
	
	log.Printf("Configuration Loaded: PORT=%s, DB_PATH=%s", port, dbPath)

	return &Config{
		Port:          port,
		DBPath:        dbPath,
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
