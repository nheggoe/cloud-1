package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logging wraps an HTTP handler to log request details,
// including method, path, status, duration, and metadata.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqID := RequestID(r)
			if reqID != "" {
				w.Header().Set("X-Request-ID", reqID)
			}

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			status := rw.status
			if status == 0 {
				status = http.StatusOK
			}

			slog.Info(
				"request completed",
				"status", status,
				"duration_ms", time.Since(start).Milliseconds(),
				"bytes", rw.bytes,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		},
	)
}

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(p)
	rw.bytes += n
	return n, err
}
