package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	countriesUpstreamPath    = "alpha/"
	countriesUpstreamTimeout = 5 * time.Second
)

// Country represents the upstream REST Countries API response shape.
type Country struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Capital    []string          `json:"capital"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Area       float64           `json:"area"`
	Population int               `json:"population"`
	Continents []string          `json:"continents"`
	Flags      struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
		Alt string `json:"alt"`
	} `json:"flags"`
}

// CountriesClient handles HTTP communication with the REST Countries API.
type CountriesClient struct {
	client  *http.Client
	baseURL string
}

// NewCountriesClient creates a CountriesClient for the given base URL.
// The base URL should include the version path, e.g. "http://129.241.150.113:8080/v3.1".
func NewCountriesClient(baseURL string) *CountriesClient {
	cleaned := strings.TrimSpace(baseURL)
	if cleaned != "" {
		cleaned = strings.TrimRight(cleaned, "/") + "/"
	}
	return &CountriesClient{
		client:  &http.Client{Timeout: countriesUpstreamTimeout},
		baseURL: cleaned,
	}
}

// GetByAlpha fetches country information by a two-letter country code.
func (c *CountriesClient) GetByAlpha(ctx context.Context, countryCode string) ([]Country, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("countries endpoint is not configured")
	}

	url := c.baseURL + countriesUpstreamPath + countryCode

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build upstream request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach countries endpoint: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("countries endpoint returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from countries endpoint: %w", err)
	}

	var countries []Country
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return countries, nil
}
