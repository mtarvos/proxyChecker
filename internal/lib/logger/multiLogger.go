package logger

import (
	"context"
	"log/slog"
)

type multiLogger struct {
	handlers []slog.Handler
}

func newMultiLogger(handlers ...slog.Handler) *multiLogger {
	return &multiLogger{handlers: handlers}
}

func (m multiLogger) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m.handlers {
		if !handler.Enabled(ctx, level) {
			return false
		}
	}

	return true
}

func (m multiLogger) Handle(ctx context.Context, record slog.Record) error {
	var err error
	for _, handler := range m.handlers {
		if handler.Enabled(ctx, record.Level) {
			err = handler.Handle(ctx, record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m multiLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))

	for i, handler := range m.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}

	return &multiLogger{handlers: handlers}
}

func (m multiLogger) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))

	for i, handler := range m.handlers {
		handlers[i] = handler.WithGroup(name)
	}

	return &multiLogger{handlers: handlers}
}
