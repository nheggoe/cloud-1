package router

import (
	"countryinfo/internal/config"
	"countryinfo/internal/handler/exchange"
	"countryinfo/internal/handler/info"
	"countryinfo/internal/handler/status"
	"countryinfo/internal/restclient"
	"net/http"
)

func New(cfg *config.Config) http.Handler {
	countriesClient := restclient.NewCountriesClient(cfg.CountriesEndpoint)
	currencyClient := restclient.NewCurrencyClient(cfg.CurrencyEndpoint)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /countryinfo/v1/status", status.Handler(cfg))
	mux.HandleFunc("GET /countryinfo/v1/info/{country_code}", info.Handler(countriesClient))
	mux.HandleFunc("GET /countryinfo/v1/exchange/{country_code}", exchange.Handler(countriesClient, currencyClient))
	return mux
}
