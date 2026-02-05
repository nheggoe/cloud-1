package endpoint

import (
	"fmt"
	"net/http"
)

const (
	InfoBase = "/countryinfo/v1/info"
	Info     = InfoBase + "/"
)

func InfoMux(endpoint string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", infoUsageHandler)
	mux.HandleFunc("GET /{country_code}", infoHandler)
	return mux
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		param := r.PathValue("two_letter_country_code")
		w.Write([]byte(param))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}

func infoUsageHandler(w http.ResponseWriter, r *http.Request) {
	usage := `The correct usage of %s is:
POST %s{two_letter_country_code}`
	http.Error(w, fmt.Sprintf(usage, Info, Info), http.StatusBadRequest)
}
