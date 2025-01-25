package middleware

import (
	"log/slog"
	"net/http"
)

type FuncMiddleware func(http.Handler) http.Handler

type Middleware struct {
	log *slog.Logger
}

func NewMiddleware(log *slog.Logger) *Middleware {
	return &Middleware{log: log}
}

func (m *Middleware) GetStack() FuncMiddleware {
	stack := m.CreateStack(
		m.Recoverer,
		m.Logging,
		m.RequestID,
	)
	return stack
}

func (m *Middleware) CreateStack(ms ...FuncMiddleware) FuncMiddleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range ms {
			next = middleware(next)
		}
		return next
	}
}
