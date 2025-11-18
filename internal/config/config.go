package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	Port       string
}

var AppConfig *Config

func Load() error {
	_ = godotenv.Load()

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "validation_user"),
		DBPassword: getEnv("DB_PASSWORD", "validation_pass"),
		DBName:     getEnv("DB_NAME", "validation_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:       getEnv("PORT", "8080"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type ServerMetrics struct {
	StartTime    time.Time
	RequestCount int64
}

var Metrics = &ServerMetrics{
	StartTime:    time.Now(),
	RequestCount: 0,
}

func (m *ServerMetrics) IncrementRequestCount() {
	m.RequestCount++
}

func (m *ServerMetrics) GetUptime() time.Duration {
	return time.Since(m.StartTime)
}

func (m *ServerMetrics) GetRequestCount() int64 {
	return m.RequestCount
}

