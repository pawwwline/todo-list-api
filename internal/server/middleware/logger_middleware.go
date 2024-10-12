package middleware

import (
	"net/http"
	"time"
	l "todo-list-api/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(wrapped, r)

		l.Logger.Info("Request received",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
		)

		l.Logger.Info("Response sent",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
		)
	})
}
