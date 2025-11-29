package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anish-chanda/ferna/internal/db"
	"github.com/anish-chanda/ferna/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret      string // JWT secret key
	TokenDuration  int    // Token duration in minutes
	CookieDuration int    // Cookie duration in hours
	BaseURL        string // Base URL for the application
	AvatarPath     string // Path for avatar storage, default to /var/lib/ferna/avatars
	DisableXSRF    bool   // Whether to disable XSRF protection, this default to true
}

// Config holds all application configuration
type Config struct {
	// Server configuration
	APIPort int
	Host    string

	// Database configuration
	Database db.Config

	// Logger configuration
	Logger logger.Config

	// Authentication configuration
	Auth AuthConfig
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		// Default server configuration
		APIPort: getEnvAsInt("API_PORT", 8080),
		Host:    getEnv("HOST", "0.0.0.0"),

		// Database configuration
		Database: db.Config{
			DSN:             getEnv("DATABASE_URL", ""),
			MaxConns:        getEnvAsInt32("DB_MAX_CONNS", 30),
			MinConns:        getEnvAsInt32("DB_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", time.Minute*30),
		},

		// Default logger configuration
		Logger: logger.Config{
			Level:   getEnv("LOG_LEVEL", "info"),
			Pretty:  getEnvAsBool("LOG_PRETTY", true),
			Service: "ferna-api",
		},

		// Authentication configuration
		Auth: AuthConfig{
			JWTSecret:      getEnv("JWT_SECRET", "ferna-jwt-secret-key-change-in-production"),
			TokenDuration:  60,   // 60 minutes
			CookieDuration: 1440, // 60 days (60 * 24 hours)
			BaseURL:        getEnv("API_BASE_URL", "http://localhost:8080"),
			AvatarPath:     "./data/avatars",
			DisableXSRF:    true,
		},
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate ensures the configuration is valid
func (c *Config) validate() error {
	if c.Database.DSN == "" {
		return errors.New("DATABASE_URL cannot be empty")
	}
	// Validate DSN syntax using pgx
	if _, err := pgxpool.ParseConfig(c.Database.DSN); err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	isValidLogLevel := false
	for _, level := range validLogLevels {
		if strings.ToLower(c.Logger.Level) == level {
			isValidLogLevel = true
			break
		}
	}
	if !isValidLogLevel {
		return errors.New("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	return nil
}

// Helper functions for environment variable parsing

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as an integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// getEnvAsInt32 gets an environment variable as int32 with a fallback value
func getEnvAsInt32(key string, fallback int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(intVal)
		}
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as duration with a fallback value
// Expects duration in format like "1h", "30m", "5s"
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as a boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "TRUE":
			return true
		case "false", "FALSE":
			return false
		}
	}
	return fallback
}
