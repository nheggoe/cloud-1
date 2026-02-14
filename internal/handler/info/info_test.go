package info

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoHandlerUsesConfiguredEndpoint(t *testing.T) {
	t.Parallel()

	gotPath := ""
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"country":"Norway"}`))
	}))
	defer upstream.Close()

	handler := Handler(upstream.URL + "/v3.1")
	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/info/no", nil)
	req.SetPathValue("country_code", "no")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if gotPath != "/v3.1/alpha/no" {
		t.Fatalf("expected upstream path /v3.1/alpha/no, got %q", gotPath)
	}
	if got := w.Body.String(); got != `{"country":"Norway"}` {
		t.Fatalf("unexpected response body: %q", got)
	}
}

func TestInfoHandlerRejectsInvalidCountryCode(t *testing.T) {
	t.Parallel()

	handler := Handler("http://example.com")
	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/info/nor", nil)
	req.SetPathValue("country_code", "nor")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
