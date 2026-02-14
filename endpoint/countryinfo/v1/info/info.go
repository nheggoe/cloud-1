package info

import (
	apiv1 "cloud1/endpoint/countryinfo/v1"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	Base = apiv1.Prefix + "/info"
	Path = Base + "/"

	upstreamPath    = "alpha/"
	upstreamTimeout = 5 * time.Second
)

type service struct {
	client  *http.Client
	baseURL string
}

func NewMux(endpoint string) http.Handler {
	s := &service{
		client: &http.Client{
			Timeout: upstreamTimeout,
		},
		baseURL: normalizeBaseURL(endpoint),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{country_code}", s.infoHandler)
	mux.HandleFunc("/", s.usageOrMethodNotAllowedHandler)
	return mux
}

func normalizeBaseURL(endpoint string) string {
	cleaned := strings.TrimSpace(endpoint)
	if cleaned == "" {
		return ""
	}
	return strings.TrimSuffix(cleaned, "/") + "/"
}

func isAlpha2CountryCode(code string) bool {
	if len(code) != 2 {
		return false
	}
	for i := 0; i < len(code); i++ {
		c := code[i]
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			return false
		}
	}
	return true
}

func (s *service) infoHandler(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToLower(strings.TrimSpace(r.PathValue("country_code")))
	if !isAlpha2CountryCode(countryCode) {
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

func (s *service) usageOrMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Allow", http.MethodGet)
	statusCode := http.StatusBadRequest
	if r.Method != http.MethodGet {
		statusCode = http.StatusMethodNotAllowed
	}

	http.Error(w, fmt.Sprintf("use GET %s{two_letter_country_code}", Path), statusCode)
}
