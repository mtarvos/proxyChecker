package logger

import (
	"fmt"
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

func InitLogger(env string, logFile string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal, envDev:
		colorTextHandler := tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		})

		if logFile != "" {
			fileHandler := newFileHandler(logFile, slog.LevelDebug)

			multi := newMultiLogger(
				colorTextHandler,
				fileHandler,
			)

			return slog.New(multi)
		}

		return slog.New(colorTextHandler)
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func newFileHandler(fileName string, level slog.Level) slog.Handler {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil
	}

	return slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: level,
	})
}
