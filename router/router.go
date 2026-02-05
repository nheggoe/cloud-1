package router

import (
	"cloud1/config"
	"cloud1/endpoint"
	"cloud1/fp"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var EndPoints = []string{
	endpoint.Status, endpoint.Info, endpoint.Exchange,
}

// NewRouter creates a new instance of an HTTP router
// configured with API endpoints and a root handler.
func NewRouter(cfg *config.Config) http.Handler {
	router := http.NewServeMux()
	Mount(router, endpoint.StatusBase, endpoint.StatusMux(cfg))
	Mount(router, endpoint.InfoBase, endpoint.InfoMux(cfg.Countries))
	Mount(router, endpoint.ExchangeBase, endpoint.ExchangeMux(cfg.Currency))
	router.HandleFunc("/", rootHelpOrNotFound)
	return router
}

// Mount registers an HTTP handler to a specific base path
// in the provided ServeMux with optional path stripping.
func Mount(mux *http.ServeMux, base string, h http.Handler) {
	base = strings.TrimSuffix(base, "/")
	mux.Handle(base+"/", http.StripPrefix(base, h))
	mux.HandleFunc(base, func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/"
		h.ServeHTTP(w, r2)
	})
}

// LoggingMiddleware wraps an HTTP handler to log request details,
// including method, path, status, duration, and metadata.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqID := requestID(r)
			if reqID != "" {
				w.Header().Set("X-Request-ID", reqID)
			}

			logger := slog.Default().With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			status := rw.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.Info(
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

func rootHelpOrNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf(`Available endpoints:
%s
`, fp.Reduce(EndPoints, func(s string, s2 string) string { return s + "\n" + s2 }))))
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

func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-Id"); id != "" {
		return id
	}
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return newRequestID()
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
