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
	currencyUpstreamTimeout = 5 * time.Second
)

// CurrencyResponse represents the upstream currency exchange API response.
type CurrencyResponse struct {
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

// CurrencyClient handles HTTP communication with the currency exchange API.
type CurrencyClient struct {
	client  *http.Client
	baseURL string
}

// NewCurrencyClient creates a CurrencyClient for the given base URL.
// The base URL should point to the currency service, e.g. "http://129.241.150.113:9090/currency".
func NewCurrencyClient(baseURL string) *CurrencyClient {
	cleaned := strings.TrimSpace(baseURL)
	if cleaned != "" {
		cleaned = strings.TrimRight(cleaned, "/") + "/"
	}
	return &CurrencyClient{
		client:  &http.Client{Timeout: currencyUpstreamTimeout},
		baseURL: cleaned,
	}
}

// GetExchangeRates fetches exchange rates for the given 3-letter currency code (ISO 4217).
func (c *CurrencyClient) GetExchangeRates(ctx context.Context, currencyCode string) (*CurrencyResponse, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("currency endpoint is not configured")
	}

	url := c.baseURL + strings.ToUpper(currencyCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build upstream request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach currency endpoint: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("currency endpoint returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from currency endpoint: %w", err)
	}

	var response CurrencyResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &response, nil
}
