package exchange

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeIsExplicitlyNotImplemented(t *testing.T) {
	t.Parallel()

	mux := NewMux("")
	req := httptest.NewRequest(http.MethodGet, "/nok/usd", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d", w.Code)
	}
}

func TestUsageAndMethodErrors(t *testing.T) {
	t.Parallel()

	mux := NewMux("")

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
