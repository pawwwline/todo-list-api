package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func SetupLogger(env string) *slog.Logger {
	switch env {
	case "local", "test":
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev", "prod":
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return Logger
}
