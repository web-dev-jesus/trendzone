// trendzone/config/config.go
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App        AppConfig
	MongoDB    MongoDBConfig
	SportsData SportsDataConfig
}

type AppConfig struct {
	Env      string
	Port     string
	Secret   string
	LogLevel logrus.Level
}

type MongoDBConfig struct {
	URI     string
	DBName  string
	Timeout time.Duration
}

type SportsDataConfig struct {
	APIKey  string
	BaseURL string
}

func Load() (*Config, error) {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	// Set default log level
	logLevel := logrus.InfoLevel
	if ll, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logLevel = ll
	}

	// Parse the MongoDB timeout
	timeout := 10 * time.Second
	if os.Getenv("MONGO_TIMEOUT") != "" {
		if t, err := time.ParseDuration(os.Getenv("MONGO_TIMEOUT") + "s"); err == nil {
			timeout = t
		}
	}

	return &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			Port:     getEnv("APP_PORT", "8080"),
			Secret:   getEnv("APP_SECRET", "default-secret-key"),
			LogLevel: logLevel,
		},
		MongoDB: MongoDBConfig{
			URI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DBName:  getEnv("MONGO_DB_NAME", "sportsdata_nfl"),
			Timeout: timeout,
		},
		SportsData: SportsDataConfig{
			APIKey:  getEnv("SPORTSDATA_API_KEY", ""),
			BaseURL: getEnv("SPORTSDATA_API_BASE_URL", "https://api.sportsdata.io/v3/nfl"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
