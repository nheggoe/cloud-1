package exchange

import (
	"countryinfo/internal/handler"
	"fmt"
	"net/http"
)

const (
	Base = handler.Prefix + "/exchange"
	Path = Base + "/"
)

func NewMux(string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{from}/{to}", notImplementedHandler)
	mux.HandleFunc("/", usageOrMethodNotAllowedHandler)
	return mux
}

func notImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "exchange endpoint is not implemented yet", http.StatusNotImplemented)
}

func usageOrMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Allow", http.MethodGet)
	statusCode := http.StatusBadRequest
	if r.Method != http.MethodGet {
		statusCode = http.StatusMethodNotAllowed
	}

	http.Error(w, fmt.Sprintf("use GET %s{from}/{to}", Path), statusCode)
}
