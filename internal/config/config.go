package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	APIKey      string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:   getEnv("PORT", "8080"),
		APIKey: getEnv("API_KEY", "apitest"),
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		cfg.DatabaseURL = dbURL
	} else {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "foodie")
		pass := getEnv("DB_PASSWORD", "foodie")
		name := getEnv("DB_NAME", "foodieapp")
		sslmode := getEnv("DB_SSLMODE", "disable")
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, pass, host, port, name, sslmode,
		)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
