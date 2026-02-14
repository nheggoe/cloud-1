package status

import (
	"context"
	"countryinfo/internal/config"
	"countryinfo/internal/handler"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	apiVersion         = handler.APIVersion
	statusProbeTimeout = 3 * time.Second
	countryProbePath   = "alpha/no"
	currencyProbePath  = "NOK"
)

type serviceHealth struct {
	CountryAPI  int    `json:"restcountriesapi"`
	CurrencyAPI int    `json:"currenciesapi"`
	Version     string `json:"version"`
	Uptime      int    `json:"uptime"`
}

type service struct {
	client           *http.Client
	countryProbeURL  string
	currencyProbeURL string
	startTime        time.Time
}

func Handler(cfg *config.Config) http.HandlerFunc {
	s := &service{
		client:           &http.Client{Timeout: statusProbeTimeout},
		countryProbeURL:  probeURL(cfg.Countries, countryProbePath),
		currencyProbeURL: probeURL(cfg.Currency, currencyProbePath),
		startTime:        time.Now(),
	}
	return s.statusHandler
}

func (s *service) newServiceHealth(ctx context.Context) (*serviceHealth, error) {
	countryStatusCode, countryErr := s.probe(ctx, s.countryProbeURL, "restcountriesapi")
	currencyStatusCode, currencyErr := s.probe(ctx, s.currencyProbeURL, "currenciesapi")

	return &serviceHealth{
		CountryAPI:  countryStatusCode,
		CurrencyAPI: currencyStatusCode,
		Version:     apiVersion,
		Uptime:      int(time.Since(s.startTime).Seconds()),
	}, errors.Join(countryErr, currencyErr)
}

func probeURL(base string, suffix string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	return strings.TrimSuffix(base, "/") + "/" + strings.TrimPrefix(suffix, "/")
}

func (s *service) probe(ctx context.Context, url string, name string) (int, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return 0, fmt.Errorf("%s endpoint is empty", name)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("build request for %s: %w", name, err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request %s failed: %w", name, err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return res.StatusCode, fmt.Errorf("%s returned %d", name, res.StatusCode)
	}

	return res.StatusCode, nil
}

func (s *service) statusHandler(w http.ResponseWriter, r *http.Request) {
	health, err := s.newServiceHealth(r.Context())
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusServiceUnavailable
		slog.WarnContext(
			r.Context(),
			"one or more upstream services are unhealthy",
			"error", err,
			"restcountriesapi", health.CountryAPI,
			"currenciesapi", health.CurrencyAPI,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(health); err != nil {
		slog.ErrorContext(r.Context(), "error while encoding response", "error", err)
	}
}
