package logger

import (
	"log/slog"
	"os"
)

type TypeHandler string

const (
	Text TypeHandler = "text"
	JSON TypeHandler = "json"
)

var log *slog.Logger

func Setup(typeHandlerStr, levelStr string) *slog.Logger {
	typeHandler := TypeHandler(typeHandlerStr)
	level := parseLevel(levelStr)

	options := &slog.HandlerOptions{Level: level}

	switch typeHandler {
	case Text:
		log = slog.New(slog.NewTextHandler(os.Stdout, options))
	case JSON:
		log = slog.New(slog.NewJSONHandler(os.Stdout, options))
	}

	return log
}

func GetLogger() *slog.Logger {
	return log
}

func parseLevel(levelStr string) slog.Level {
	switch levelStr {
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
