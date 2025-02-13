package mLogger

import (
	"context"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
)

type ctxKey string

const slogFields ctxKey = "slogFields"

type MLogger struct {
	slog.Handler
}

func NewMHandler(w io.Writer, opts *tint.Options) slog.Handler {
	return &MLogger{Handler: tint.NewHandler(w, opts)}
}

func (m MLogger) Handle(ctx context.Context, record slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, attr := range attrs {
			record.AddAttrs(attr)
		}
	}
	return m.Handler.Handle(ctx, record)
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}
