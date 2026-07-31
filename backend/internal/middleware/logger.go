package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			next.ServeHTTP(w, r)
			
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", time.Since(start),
			)
		})
	}
}
