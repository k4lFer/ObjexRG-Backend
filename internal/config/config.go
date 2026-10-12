package config

import (
	"os"
	"strconv"
)

type Config struct {
	BasePath	string
	Port      string
	DBUrl     string
	JWTSecret string
	JWTSecretRefresh string
	MaxUploadBytes int64
}

const defaultMaxUploadBytes int64 = 512 * 1024 * 1024 // 512 MiB

func LoadConfig() *Config {

	return &Config{
		Port:      getEnv("PORT", "8080"),
		BasePath:   os.Getenv("OBJEXRG_STORAGE_PATH"),
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_ACCESS_SECRET"),
		JWTSecretRefresh: os.Getenv("JWT_REFRESH_SECRET"),
		MaxUploadBytes: getEnvInt64("MAX_UPLOAD_BYTES", defaultMaxUploadBytes),
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt64(key string, fallback int64) int64 {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
