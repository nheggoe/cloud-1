package router

import (
	"cloud1/config"
	"cloud1/endpoint/countryinfo/v1/exchange"
	"cloud1/endpoint/countryinfo/v1/info"
	"cloud1/endpoint/countryinfo/v1/status"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Endpoints = []string{
	status.Path, info.Path, exchange.Path,
}

// NewRouter creates a new instance of an HTTP router
// configured with API endpoints and a root handler.
func NewRouter(cfg *config.Config) http.Handler {
	if cfg == nil {
		panic("config is required")
	}
	router := http.NewServeMux()
	Mount(router, status.Path, status.NewMux(cfg))
	Mount(router, info.Path, info.NewMux(cfg.Countries))
	Mount(router, exchange.Path, exchange.NewMux(cfg.Currency))
	router.HandleFunc("/", rootHelpOrNotFound)
	return router
}

func Mount(httpMux http.Handler, path string, handler http.Handler) {
	mux := httpMux.(*http.ServeMux)
	prefix := strings.TrimSuffix(path, "/")
	stripped := http.StripPrefix(prefix, handler)
	if strings.HasSuffix(path, "/") {
		mux.Handle(path, stripped)
	} else {
		// For exact paths (no trailing slash), StripPrefix yields an empty
		// URL path which won't match "/" in the sub-mux. Rewrite to "/".
		mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = "/"
			r.URL.RawPath = "/"
			handler.ServeHTTP(w, r)
		}))
		mux.Handle(path+"/", stripped)
	}
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

func rootHelpOrNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Available endpoints:\n%s\n", strings.Join(Endpoints, "\n"))
	if _, err := w.Write([]byte(body)); err != nil {
		slog.Error("failed writing root response", "error", err)
	}
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
