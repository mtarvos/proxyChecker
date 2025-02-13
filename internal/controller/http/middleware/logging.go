package middleware

import (
	"log/slog"
	"net/http"
	"proxyChecker/internal/lib/logging"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		log := logging.L(r.Context())
		log.Info(
			r.Method,
			slog.String("path", r.URL.Path),
			slog.Duration("time", time.Since(start)),
			slog.Int("status", wrapped.statusCode),
		)
	})
}
