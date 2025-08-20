package config

import (
	"fmt"
	postgres "github.com/Base-111/backend/pkg/repository"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Base-111/backend/pkg/auth"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Auth        auth.Config
	RedisConfig redis.Options
	Postgres    postgres.Config
	AppPort     string
	AppHost     string
}

func LoadFromEnv() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(currentDir, ".env")
	err = godotenv.Load(envPath)

	if err != nil {
		slog.Error("Error loading .env file", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		slog.Error("Database redis invalid : ", err)
	}
	return &Config{
		Auth: auth.Config{
			BaseURL:      os.Getenv("AUTH_BASE_URL"),
			ClientID:     os.Getenv("AUTH_CLIENT_ID"),
			RedirectURL:  os.Getenv("AUTH_REDIRECT_URL"),
			ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
			Realm:        os.Getenv("AUTH_ENVIRONMENT"),
		},
		RedisConfig: redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		Postgres: postgres.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		AppPort: os.Getenv("APP_PORT"),
		AppHost: os.Getenv("APP_HOST"),
	}, nil
}
