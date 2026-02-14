package exchange

import "net/http"

func Handler(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "exchange endpoint is not implemented yet", http.StatusNotImplemented)
	}
}
