package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port    int
	Version string
	Commit  string
	Branch  string
	BuildDate string

	// Database
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string

	// MinIO
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	// Rate Limiting
	RateLimitBuilds int // per hour per IP
	RateLimitReads  int // per hour per IP
}

func Load() *Config {
	return &Config{
		// Server
		Port:      getEnvInt("API_PORT", 8080),
		Version:   getEnv("API_VERSION", "1.0.1-pre-alpha"),
		Commit:    getEnv("API_COMMIT", "unknown"),
		Branch:    getEnv("API_BRANCH", "main"),
		BuildDate: getEnv("API_BUILD_DATE", "unknown"),

		// Database
		DBHost:     getEnv("DB_HOST", "192.168.22.137"),
		DBPort:     getEnvInt("DB_PORT", 3306),
		DBName:     getEnv("DB_NAME", "pwtoolkit"),
		DBUser:     getEnv("DB_USER", "pwtoolkit"),
		DBPassword: getEnv("DB_PASSWORD", ""),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "192.168.22.137"),
		RedisPort:     getEnvInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// MinIO
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "192.168.22.137:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinIOBucket:    getEnv("MINIO_BUCKET", "pwtoolkit"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),

		// Rate Limiting
		RateLimitBuilds: getEnvInt("RATE_LIMIT_BUILDS", 30),
		RateLimitReads:  getEnvInt("RATE_LIMIT_READS", 1000),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
