package config

import (
	"os"
)

// Config system configuration
type Config struct {
	ReddioAPIKey string
	ReddioURL    string
}

// Load loads configuration
func Load() *Config {
	return &Config{
		ReddioAPIKey: getEnv("REDDIO_API_KEY", "mock_api_key"),
		ReddioURL:    getEnv("REDDIO_URL", "https://reddio-service-prod.reddio.com"),
	}
}

// getEnv gets environment variable, returns default value if not exists
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
