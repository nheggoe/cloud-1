package router

import (
	"countryinfo/internal/config"
	"countryinfo/internal/handler/exchange"
	"countryinfo/internal/handler/info"
	"countryinfo/internal/handler/status"
	"fmt"
	"net/http"
	"strings"
)

var Endpoints = []string{
	status.Path, info.Path, exchange.Path,
}

// New creates a new instance of an HTTP router
// configured with API endpoints and a root handler.
func New(cfg *config.Config) http.Handler {
	if cfg == nil {
		panic("config is required")
	}
	router := http.NewServeMux()
	mount(router, status.Path, status.NewMux(cfg))
	mount(router, info.Path, info.NewMux(cfg.Countries))
	mount(router, exchange.Path, exchange.NewMux(cfg.Currency))
	router.HandleFunc("/", rootHelpOrNotFound)
	return router
}

func mount(httpMux http.Handler, path string, handler http.Handler) {
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

func rootHelpOrNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Available endpoints:\n%s\n", strings.Join(Endpoints, "\n"))
	if _, err := w.Write([]byte(body)); err != nil {
		// Error logging is handled by middleware
	}
}
