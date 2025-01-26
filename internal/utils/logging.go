package utils

import (
	"log/slog"
	"os"
	"strings"
)

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func LogFatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

func ConfigureLogging(level string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(level),
	}))

	slog.SetDefault(logger)
}
