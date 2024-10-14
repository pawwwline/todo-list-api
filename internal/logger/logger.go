package logger

import (
	"errors"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func SetupLogger(env string) (*slog.Logger, error) {
	if env == "" {
		return nil, errors.New("env is none")
	}
	switch env {
	case "local", "test":
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev", "prod":
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return Logger, nil
}
