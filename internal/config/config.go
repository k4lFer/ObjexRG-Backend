package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	BasePath         string
	Port             string
	DBUrl            string
	JWTSecret        string
	JWTSecretRefresh string
	MaxUploadBytes   int64
	PublicBaseURL    string

	StorageDriver    string
	S3Endpoint       string
	S3PublicEndpoint string
	S3AccessKey      string
	S3SecretKey      string
	S3Bucket         string
	S3Region         string
	S3UseSSL         bool
	S3ForcePathStyle bool
	PresignTTL       time.Duration
	SecretsMasterKey string
}

const defaultMaxUploadBytes int64 = 512 * 1024 * 1024 // 512 MiB
const defaultPresignTTL = 15 * time.Minute

func LoadConfig() *Config {

	return &Config{
		Port:             getEnv("PORT", "8080"),
		BasePath:         os.Getenv("OBJEXRG_STORAGE_PATH"),
		DBUrl:            os.Getenv("DATABASE_URL"),
		JWTSecret:        os.Getenv("JWT_ACCESS_SECRET"),
		JWTSecretRefresh: os.Getenv("JWT_REFRESH_SECRET"),
		MaxUploadBytes:   getEnvInt64("MAX_UPLOAD_BYTES", defaultMaxUploadBytes),
		PublicBaseURL:    getEnv("PUBLIC_BASE_URL", "http://localhost:7080"),

		StorageDriver:    getEnv("STORAGE_DRIVER", "local"),
		S3Endpoint:       os.Getenv("S3_ENDPOINT"),
		S3PublicEndpoint: os.Getenv("S3_PUBLIC_ENDPOINT"),
		S3AccessKey:      os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:      os.Getenv("S3_SECRET_KEY"),
		S3Bucket:         getEnv("S3_BUCKET", "objexrg"),
		S3Region:         getEnv("S3_REGION", "us-east-1"),
		S3UseSSL:         getEnvBool("S3_USE_SSL", false),
		S3ForcePathStyle: getEnvBool("S3_FORCE_PATH_STYLE", true),
		PresignTTL:       getEnvDuration("PRESIGN_TTL", defaultPresignTTL),
		SecretsMasterKey: os.Getenv("SECRETS_MASTER_KEY"),
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

func getEnvBool(key string, fallback bool) bool {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
