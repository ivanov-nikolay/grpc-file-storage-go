package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort      string
	Database      DatabaseConfig
	StoragePath   string
	UploadLimit   int64
	DownloadLimit int64
	ListLimit     int64
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string

	CreateDB bool

	SuperUser     string
	SuperPassword string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		Database: DatabaseConfig{
			Host:          getEnv("DATABASE_HOST", "localhost"),
			Port:          getEnv("DATABASE_PORT", "5432"),
			User:          getEnv("DATABASE_USER", "postgres"),
			Password:      getEnv("DATABASE_PASSWORD", "12345"),
			DBName:        getEnv("DATABASE_NAME", "file_storage"),
			SSLMode:       getEnv("DATABASE_SSLMODE", "disable"),
			CreateDB:      getEnvBool("DATABASE_CREATE_DB", true),
			SuperUser:     getEnv("DATABASE_SUPER_USER", "postgres"),
			SuperPassword: getEnv("DATABASE_SUPER_PASSWORD", "12345"),
		},
		StoragePath:   getEnv("STORAGE_PATH", "./storage/files"),
		UploadLimit:   getEnvInt64("UPLOAD_LIMIT", 10),
		DownloadLimit: getEnvInt64("DOWNLOAD_LIMIT", 10),
		ListLimit:     getEnvInt64("LIST_LIMIT", 100),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}

	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}

	return defaultValue
}
