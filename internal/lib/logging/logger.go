package logging

import (
	"context"
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
	switch env {
	case envLocal, envDev:
		return NewLogger(slog.LevelDebug, false, false)
	default:
		return NewLogger(slog.LevelInfo, true, true)
	}
}

func NewLogger(level slog.Level, addSource bool, isJSON bool) *slog.Logger {
	var logger *slog.Logger
	if isJSON {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: addSource,
		}))
	} else {
		logger = slog.New(tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
			Level:      level,
			AddSource:  addSource,
			TimeFormat: time.DateTime,
		}))
	}

	slog.SetDefault(logger)

	return logger
}

func L(ctx context.Context) *slog.Logger {
	return loggerFromContext(ctx)
}

func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "nil")
	}

	return slog.String("error", err.Error())
}
