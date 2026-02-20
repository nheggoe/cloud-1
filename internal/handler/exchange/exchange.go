package exchange

import (
	"countryinfo/internal/restclient"
	"countryinfo/internal/util"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"
)

type ExchangeResponse struct {
	Country       string               `json:"country"`
	BaseCurrency  string               `json:"base-currency"`
	ExchangeRates []map[string]float64 `json:"exchange-rates"`
}

type service struct {
	countries  *restclient.CountriesClient
	currencies *restclient.CurrencyClient
}

func Handler(countries *restclient.CountriesClient, currencies *restclient.CurrencyClient) http.HandlerFunc {
	s := &service{
		countries:  countries,
		currencies: currencies,
	}
	return s.exchangeHandler
}

func (s *service) exchangeHandler(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToLower(strings.TrimSpace(r.PathValue("country_code")))
	if !util.IsTwoLetterCountryCode(countryCode) {
		http.Error(
			w,
			fmt.Sprintf("%s\ninvalid country code: %s", http.StatusText(http.StatusBadRequest), countryCode),
			http.StatusBadRequest,
		)
		return
	}

	ctx := r.Context()

	// 1. Look up the input country to get its currency and borders.
	countries, err := s.countries.GetByAlpha(ctx, countryCode)
	if err != nil {
		slog.ErrorContext(ctx, "failed to look up country", "error", err, "country_code", countryCode)
		http.Error(w, "failed to look up country", http.StatusBadGateway)
		return
	}
	if len(countries) == 0 {
		http.Error(w, "country not found", http.StatusNotFound)
		return
	}
	country := countries[0]

	// Extract the base currency code (first key in the currencies map).
	baseCurrencyCode := firstCurrencyCode(country)
	if baseCurrencyCode == "" {
		http.Error(w, "no currency found for country", http.StatusNotFound)
		return
	}

	if len(country.Borders) == 0 {
		// Country has no land borders, return empty exchange rates.
		writeJSON(w, r, ExchangeResponse{
			Country:       country.Name.Common,
			BaseCurrency:  baseCurrencyCode,
			ExchangeRates: []map[string]float64{},
		})
		return
	}

	// 2. Look up each bordering country to collect their currency codes.
	neighbourCurrencies := make(map[string]struct{})
	for _, borderCode := range country.Borders {
		neighbours, err := s.countries.GetByAlpha(ctx, strings.ToLower(borderCode))
		if err != nil {
			slog.WarnContext(ctx, "failed to look up border country", "error", err, "border_code", borderCode)
			continue
		}
		if len(neighbours) == 0 {
			continue
		}
		for code := range neighbours[0].Currencies {
			neighbourCurrencies[code] = struct{}{}
		}
	}

	// 3. Fetch exchange rates for the base currency.
	rates, err := s.currencies.GetExchangeRates(ctx, baseCurrencyCode)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch exchange rates", "error", err, "base_currency", baseCurrencyCode)
		http.Error(w, "failed to fetch exchange rates", http.StatusBadGateway)
		return
	}

	// 4. Filter rates to only include neighboring countries' currencies.
	codes := make([]string, 0, len(neighbourCurrencies))
	for code := range neighbourCurrencies {
		if _, ok := rates.Rates[code]; ok {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)

	exchangeRates := make([]map[string]float64, 0, len(codes))
	for _, code := range codes {
		exchangeRates = append(exchangeRates, map[string]float64{code: rates.Rates[code]})
	}

	writeJSON(w, r, ExchangeResponse{
		Country:       country.Name.Common,
		BaseCurrency:  baseCurrencyCode,
		ExchangeRates: exchangeRates,
	})

	slog.InfoContext(ctx, "exchange request completed",
		"country_code", countryCode,
		"base_currency", baseCurrencyCode,
		"neighbour_currencies", len(exchangeRates),
	)
}

func firstCurrencyCode(c restclient.Country) string {
	for code := range c.Currencies {
		return code
	}
	return ""
}

func writeJSON(w http.ResponseWriter, r *http.Request, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal json", "error", err)
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
