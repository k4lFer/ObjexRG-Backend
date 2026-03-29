package config

import (
	"os"
)

type Config struct {
	BasePath	string
	Port      string
	DBUrl     string
	JWTSecret string
	JWTSecretRefresh string
}

func LoadConfig() *Config {

	return &Config{
		Port:      getEnv("PORT", "8080"),
		BasePath:   os.Getenv("OBJEXRG_STORAGE_PATH"),
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_ACCESS_SECRET"),
		JWTSecretRefresh: os.Getenv("JWT_REFRESH_SECRET"),
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}