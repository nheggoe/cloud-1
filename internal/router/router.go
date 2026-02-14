package router

import (
	"countryinfo/internal/config"
	"countryinfo/internal/handler/exchange"
	"countryinfo/internal/handler/info"
	"countryinfo/internal/handler/status"
	"net/http"
)

func New(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /countryinfo/v1/status", status.Handler(cfg))
	mux.HandleFunc("GET /countryinfo/v1/info/{country_code}", info.Handler(cfg.Countries))
	mux.HandleFunc("GET /countryinfo/v1/exchange/{country_code}", exchange.Handler(cfg.Currency))
	return mux
}
