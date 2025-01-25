package logger

import (
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func InitLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal, envDev:
		logger = slog.New(
			tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.DateTime,
			}))
		return logger
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
