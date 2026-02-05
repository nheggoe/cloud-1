package endpoint

import (
	"cloud1/config"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

const (
	StatusBase = "/countryinfo/v1/status"
	Status     = StatusBase + "/"
)

var startTime = time.Now()

const apiVersion = "v1"

type serviceHealth struct {
	CountryAPI  string `json:"restcountriesapi"`
	CurrencyAPI string `json:"currenciesapi"`
	Version     string `json:"version"`
	Uptime      int    `json:"uptime"`
}

func StatusMux(ctf *config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", statusHandler)
	mux.HandleFunc("/", statusUsageHandler)
	return mux
}

func newServiceHealth() *serviceHealth {
	return &serviceHealth{
		CountryAPI:  "",
		CurrencyAPI: "",
		Version:     apiVersion,
		Uptime:      int(time.Since(startTime).Seconds()),
	}

}

func statusUsageHandler(w http.ResponseWriter, r *http.Request) {

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(newServiceHealth())
		if err != nil {
			slog.ErrorContext(r.Context(), "Error while encoding response", "error", err)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
