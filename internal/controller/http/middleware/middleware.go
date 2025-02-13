package middleware

import (
	"net/http"
)

type FuncMiddleware func(http.Handler) http.Handler

type Middleware struct {
	Env string
}

func NewMiddleware(env string) *Middleware {
	return &Middleware{Env: env}
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
