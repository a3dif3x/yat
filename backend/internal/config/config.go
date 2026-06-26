package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Port        string
	LogLevel    slog.Level
	DatabaseURL string
}

func Load() (Config, error) {
	port := getEnv("PORT", "8080")
	logLevel, err := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		return Config{}, err
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("database URL is required")
	}

	cfg := Config{
		Port:        port,
		LogLevel:    logLevel,
		DatabaseURL: databaseURL,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func parseLogLevel(s string) (slog.Level, error) {
	switch s {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid LOG_LEVEL %q (want debug, info, warn, or error)", s)
	}
}
