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

	mux := NewMux(upstream.URL + "/v3.1")
	req := httptest.NewRequest(http.MethodGet, "/no", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

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

	mux := NewMux("http://example.com")
	req := httptest.NewRequest(http.MethodGet, "/nor", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUsageAndMethodErrors(t *testing.T) {
	t.Parallel()

	mux := NewMux("http://example.com")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for GET /, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for POST /, got %d", w.Code)
	}
}
