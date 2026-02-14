package router

import (
	"countryinfo/internal/config"
	"countryinfo/internal/handler/exchange"
	"countryinfo/internal/handler/info"
	"countryinfo/internal/handler/status"
	"fmt"
	"net/http"
)

func New(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Method: %s\nPath: %s\n", r.Method, r.URL.Path)
	})
	mux.HandleFunc("GET /countryinfo/v1/status", status.Handler(cfg))
	mux.HandleFunc("GET /countryinfo/v1/info/{country_code}", info.Handler(cfg.Countries))
	mux.HandleFunc("GET /countryinfo/v1/exchange/{country_code}", exchange.Handler(cfg.Currency))
	return mux
}
