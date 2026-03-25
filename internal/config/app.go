package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv          string
	ServerPort      string
	DBDSN           string
	JWTSecrect      string
	JWTExpiresHours int
	RedisHost       string
	RedisPort       string
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ENV must have %s ", key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	cfg := AppConfig{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: mustGetEnv("SERVER_PORT"),
		DBDSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mustGetEnv("DB_USER"),
			mustGetEnv("DB_PASSWORD"),
			mustGetEnv("DB_HOST"),
			mustGetEnv("DB_PORT"),
			mustGetEnv("DB_NAME"),
		),
		JWTSecrect: mustGetEnv("JWT_SECRECT_KEY"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
	}

	hours, err := strconv.Atoi(mustGetEnv("JWT_EXPIRES_IN"))
	if err != nil || hours <= 0 {
		log.Fatal("JWT_EXPIRES_HOURS must be positive")
	}
	cfg.JWTExpiresHours = hours
	return cfg
}
