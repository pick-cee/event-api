package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	RedisURL   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:       GetEnv("PORT", ""),
		DBHost:     GetEnv("DB_HOST", ""),
		DBPort:     GetEnv("DB_PORT", ""),
		DBUser:     GetEnv("DB_USER", ""),
		DBPassword: GetEnv("DB_PASSWORD", ""),
		DBName:     GetEnv("DB_NAME", ""),
		JWTSecret:  GetEnv("JWT_SECRET", ""),
		RedisURL:   GetEnv("REDIS_URL", ""),
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
