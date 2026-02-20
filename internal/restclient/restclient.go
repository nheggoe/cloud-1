package restclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Probe performs a lightweight HTTP GET against url and returns the status code.
// It is used by the status handler to check upstream API health.
func Probe(ctx context.Context, client *http.Client, url string, name string) (int, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return 0, fmt.Errorf("%s endpoint is empty", name)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("build request for %s: %w", name, err)
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request %s failed: %w", name, err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return res.StatusCode, fmt.Errorf("%s returned %d", name, res.StatusCode)
	}

	return res.StatusCode, nil
}
