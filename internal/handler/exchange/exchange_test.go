package exchange

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeIsExplicitlyNotImplemented(t *testing.T) {
	t.Parallel()

	handler := Handler("")
	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/exchange/nok/usd", nil)
	req.SetPathValue("from", "nok")
	req.SetPathValue("to", "usd")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d", w.Code)
	}
}
