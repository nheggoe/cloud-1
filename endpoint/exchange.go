package endpoint

import "net/http"

const (
	ExchangeBase = "/countryinfo/v1/exchange"
	Exchange     = ExchangeBase + "/"
)

func ExchangeMux(string) http.Handler {
	mux := http.NewServeMux()
	return mux
}
