package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWrite struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWrite) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sTime := time.Now()
			logger.InfoContext(
				context.Background(),
				"Incoming HTTP request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
			lrw := &loggingResponseWrite{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)
			duration := time.Since(sTime)
			logger.InfoContext(
				context.Background(),
				"Completed HTTP request",
				slog.Int("status", lrw.statusCode),
				slog.Duration("duration", duration),
			)
		})
	}
}
