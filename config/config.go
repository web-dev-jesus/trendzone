package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// MongoDB Configuration
	MongoURI      string
	MongoDatabase string

	// SportsData.io API Configuration
	SportsDataAPIKey  string
	SportsDataBaseURL string

	// Application Configuration
	Season     int
	SeasonType string
	LogLevel   string

	// Security Configuration
	EnableTLS   bool
	TLSCertFile string
	TLSKeyFile  string

	// Rate Limiting
	APICallDelay time.Duration
}

// LoadConfig loads application configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		// MongoDB Configuration - Required
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: getEnvWithDefault("MONGO_DATABASE", "sportsdata_nfl"),

		// SportsData.io API Configuration - Required
		SportsDataAPIKey:  os.Getenv("SPORTSDATA_API_KEY"),
		SportsDataBaseURL: getEnvWithDefault("SPORTSDATA_BASE_URL", "https://api.sportsdata.io/v3/nfl"),

		// Application Configuration
		Season:     getEnvAsIntWithDefault("SEASON", 2023),
		SeasonType: getEnvWithDefault("SEASON_TYPE", "REG"),
		LogLevel:   getEnvWithDefault("LOG_LEVEL", "info"),

		// Security Configuration
		EnableTLS:   getEnvAsBoolWithDefault("ENABLE_TLS", false),
		TLSCertFile: getEnvWithDefault("TLS_CERT_FILE", ""),
		TLSKeyFile:  getEnvWithDefault("TLS_KEY_FILE", ""),

		// Rate Limiting
		APICallDelay: time.Duration(getEnvAsIntWithDefault("API_CALL_DELAY_MS", 1000)) * time.Millisecond,
	}

	// Validate required configuration
	if config.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}

	if config.SportsDataAPIKey == "" {
		return nil, fmt.Errorf("SPORTSDATA_API_KEY is required")
	}

	return config, nil
}

// Helper functions to read environment variables with default values
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolWithDefault(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
