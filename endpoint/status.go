package endpoint

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

const apiVersion = "v1"

type serviceHealth struct {
	CountryAPI  string `json:"restcountriesapi"`
	CurrencyAPI string `json:"currenciesapi"`
	Version     string `json:"version"`
	Uptime      int    `json:"uptime"`
}

func newServiceHealth() *serviceHealth {
	return &serviceHealth{
		CountryAPI:  "",
		CurrencyAPI: "",
		Version:     apiVersion,
		Uptime:      int(time.Since(startTime).Seconds()),
	}

}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		res, err := json.Marshal(newServiceHealth())
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(append(res, '\n'))
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
