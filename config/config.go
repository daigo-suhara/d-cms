package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBURL        string
	R2BucketName string
	R2Endpoint   string
	R2PublicBase string
	AWSKeyID     string
	AWSKeySecret string
	AWSRegion    string
	AdminToken   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DBURL:        mustGetEnv("DB_URL"),
		R2BucketName: getEnv("R2_BUCKET_NAME", ""),
		R2Endpoint:   getEnv("R2_ENDPOINT", ""),
		R2PublicBase: getEnv("R2_PUBLIC_BASE_URL", ""),
		AWSKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSKeySecret: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:    getEnv("AWS_REGION", "auto"),
		AdminToken:   mustGetEnv("ADMIN_TOKEN"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}
