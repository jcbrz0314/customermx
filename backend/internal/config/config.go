package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Provider     string
	From         string
	FrontendURL  string
	LogoURL      string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPUseTLS   bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "customermx"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessTTL:     parseDuration(getEnv("JWT_ACCESS_TTL", "15m"), 15*time.Minute),
			RefreshTTL:    parseDuration(getEnv("JWT_REFRESH_TTL", "168h"), 168*time.Hour),
		},
		Email: EmailConfig{
			Provider:     getEnv("EMAIL_PROVIDER", "mock"),
			From:         getEnv("EMAIL_FROM", "noreply@customermx.com"),
			FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
			LogoURL:      getEnv("EMAIL_LOGO_URL", ""),
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     parseInt(getEnv("SMTP_PORT", "587"), 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPUseTLS:   parseBool(getEnv("SMTP_USE_TLS", "true")),
		},
	}

	// Validate required fields
	if config.JWT.AccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if config.JWT.RefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string or returns a default value
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return defaultValue
}

// parseBool parses a boolean string or returns false
func parseBool(value string) bool {
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	return false
}

// parseInt parses an integer string or returns the default value
func parseInt(value string, defaultValue int) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return defaultValue
}
