package exchange

import (
	"countryinfo/internal/restclient"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeHandlerRejectsInvalidCountryCode(t *testing.T) {
	t.Parallel()

	countries := restclient.NewCountriesClient("http://example.com")
	currencies := restclient.NewCurrencyClient("http://example.com")
	handler := Handler(countries, currencies)

	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/exchange/nor", nil)
	req.SetPathValue("country_code", "nor")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestExchangeHandlerReturnsNeighbourRates(t *testing.T) {
	t.Parallel()

	// Mock countries API: serves Norway and its neighbours.
	countriesAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := map[string]string{
			"/v3.1/alpha/no":  `[{"name":{"common":"Norway"},"currencies":{"NOK":{"name":"Norwegian Krone","symbol":"kr"}},"borders":["FIN","SWE"],"capital":["Oslo"],"continents":["Europe"],"population":5379475,"area":323802,"languages":{"nob":"Norwegian Bokmal"},"flags":{"png":"","svg":"","alt":""}}]`,
			"/v3.1/alpha/fin": `[{"name":{"common":"Finland"},"currencies":{"EUR":{"name":"Euro","symbol":"€"}},"borders":["NOR","SWE","RUS"],"capital":["Helsinki"],"continents":["Europe"],"population":5530719,"area":338424,"languages":{"fin":"Finnish"},"flags":{"png":"","svg":"","alt":""}}]`,
			"/v3.1/alpha/swe": `[{"name":{"common":"Sweden"},"currencies":{"SEK":{"name":"Swedish Krona","symbol":"kr"}},"borders":["NOR","FIN","DNK"],"capital":["Stockholm"],"continents":["Europe"],"population":10353442,"area":450295,"languages":{"swe":"Swedish"},"flags":{"png":"","svg":"","alt":""}}]`,
		}

		w.Header().Set("Content-Type", "application/json")
		if body, ok := responses[r.URL.Path]; ok {
			_, _ = w.Write([]byte(body))
		} else {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"status":404,"message":"Not Found"}`))
		}
	}))
	defer countriesAPI.Close()

	// Mock currency API: serves NOK exchange rates.
	currencyAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"success","base_code":"NOK","rates":{"NOK":1,"EUR":0.086536,"SEK":0.914075,"USD":0.105048}}`))
	}))
	defer currencyAPI.Close()

	countries := restclient.NewCountriesClient(countriesAPI.URL + "/v3.1")
	currencies := restclient.NewCurrencyClient(currencyAPI.URL + "/currency")
	handler := Handler(countries, currencies)

	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/exchange/no", nil)
	req.SetPathValue("country_code", "no")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp ExchangeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Country != "Norway" {
		t.Errorf("expected country Norway, got %q", resp.Country)
	}
	if resp.BaseCurrency != "NOK" {
		t.Errorf("expected base currency NOK, got %q", resp.BaseCurrency)
	}
	if len(resp.ExchangeRates) != 2 {
		t.Fatalf("expected 2 exchange rates (EUR, SEK), got %d: %v", len(resp.ExchangeRates), resp.ExchangeRates)
	}

	// Collect the rates into a flat map for easier assertion.
	got := make(map[string]float64)
	for _, entry := range resp.ExchangeRates {
		for k, v := range entry {
			got[k] = v
		}
	}
	if got["EUR"] != 0.086536 {
		t.Errorf("expected EUR rate 0.086536, got %f", got["EUR"])
	}
	if got["SEK"] != 0.914075 {
		t.Errorf("expected SEK rate 0.914075, got %f", got["SEK"])
	}
	if _, ok := got["USD"]; ok {
		t.Error("USD should not be in exchange rates (not a neighbour currency)")
	}
}

func TestExchangeHandlerCountryWithNoBorders(t *testing.T) {
	t.Parallel()

	countriesAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":{"common":"Iceland"},"currencies":{"ISK":{"name":"Icelandic Krona","symbol":"kr"}},"borders":[],"capital":["Reykjavik"],"continents":["Europe"],"population":366425,"area":103000,"languages":{"isl":"Icelandic"},"flags":{"png":"","svg":"","alt":""}}]`))
	}))
	defer countriesAPI.Close()

	currencyAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("currency API should not be called for a country with no borders")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer currencyAPI.Close()

	countries := restclient.NewCountriesClient(countriesAPI.URL + "/v3.1")
	currencies := restclient.NewCurrencyClient(currencyAPI.URL + "/currency")
	handler := Handler(countries, currencies)

	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/exchange/is", nil)
	req.SetPathValue("country_code", "is")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp ExchangeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Country != "Iceland" {
		t.Errorf("expected country Iceland, got %q", resp.Country)
	}
	if len(resp.ExchangeRates) != 0 {
		t.Errorf("expected empty exchange rates, got %v", resp.ExchangeRates)
	}
}
