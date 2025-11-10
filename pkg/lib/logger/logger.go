package logger

import (
	"log/slog"
	"os"
)

const (
	Local string = "local"
	Dev   string = "dev"
	Prod  string = "prod"
)

var log *slog.Logger

func Setup(env string) *slog.Logger {
	switch env {
	case Local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func GetLogger() *slog.Logger {
	return log
}
