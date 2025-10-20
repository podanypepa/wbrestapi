// Package config provides application configuration
package config

import (
	"cmp"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	RateLimitMax    int
	RateLimitWindow time.Duration
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host            string
	User            string
	Password        string
	Name            string
	Port            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            cmp.Or(os.Getenv("PORT"), "3000"),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 5*time.Second),
			RateLimitMax:    getIntEnv("RATE_LIMIT_MAX", 100),
			RateLimitWindow: getDurationEnv("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		Database: DatabaseConfig{
			Host:            cmp.Or(os.Getenv("DB_HOST"), "localhost"),
			User:            cmp.Or(os.Getenv("DB_USER"), "postgres"),
			Password:        os.Getenv("DB_PASSWORD"),
			Name:            cmp.Or(os.Getenv("DB_NAME"), "users"),
			Port:            cmp.Or(os.Getenv("DB_PORT"), "5432"),
			SSLMode:         cmp.Or(os.Getenv("DB_SSL"), "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
	}
}

func getIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
