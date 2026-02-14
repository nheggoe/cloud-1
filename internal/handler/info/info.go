package info

import (
	"countryinfo/internal/util"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	upstreamPath    = "alpha/"
	upstreamTimeout = 5 * time.Second
)

type service struct {
	client  *http.Client
	baseURL string
}

func Handler(endpoint string) http.HandlerFunc {
	s := &service{
		client: &http.Client{
			Timeout: upstreamTimeout,
		},
		baseURL: util.CleanUrl(endpoint),
	}
	return s.infoHandler
}

func (s *service) infoHandler(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToLower(strings.TrimSpace(r.PathValue("country_code")))
	if !util.IsTwoLetterCountryCode(countryCode) {
		http.Error(
			w,
			fmt.Sprintf("%s\ninvalid country code: %s", http.StatusText(http.StatusBadRequest), countryCode),
			http.StatusBadRequest,
		)
		return
	}
	if s.baseURL == "" {
		http.Error(w, "countries endpoint is not configured", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, s.baseURL+upstreamPath+countryCode, nil)
	if err != nil {
		http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
		return
	}

	res, err := s.client.Do(req)
	if err != nil {
		http.Error(w, "failed to reach countries endpoint", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "countries endpoint returned non-JSON response", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)

	if _, err := io.Copy(w, res.Body); err != nil {
		slog.ErrorContext(r.Context(), "failed to proxy response body", "error", err)
		return
	}

	slog.InfoContext(
		r.Context(),
		"country info request completed",
		"country_code", countryCode,
		"status", res.StatusCode,
	)
}
