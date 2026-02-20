package info

import (
	"countryinfo/internal/restclient"
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
		_, _ = w.Write([]byte(`[{"name":{"common":"Norway"},"capital":["Oslo"],"continents":["Europe"],"population":5379475,"area":323802,"languages":{"nno":"Norwegian Nynorsk","nob":"Norwegian Bokmal","smi":"Sami"},"borders":["FIN","SWE","RUS"],"flags":{"png":"https://flagcdn.com/w320/no.png","svg":"https://flagcdn.com/no.svg","alt":"Norway flag"}}]`))
	}))
	defer upstream.Close()

	client := restclient.NewCountriesClient(upstream.URL + "/v3.1")
	handler := Handler(client)
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

	expected := `{"name":"Norway","continents":["Europe"],"population":5379475,"area":323802,"languages":{"nno":"Norwegian Nynorsk","nob":"Norwegian Bokmal","smi":"Sami"},"borders":["FIN","SWE","RUS"],"flag":"https://flagcdn.com/w320/no.png","capital":"Oslo"}`
	if got := w.Body.String(); got != expected {
		t.Fatalf("unexpected response body:\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestInfoHandlerRejectsInvalidCountryCode(t *testing.T) {
	t.Parallel()

	client := restclient.NewCountriesClient("http://example.com")
	handler := Handler(client)
	req := httptest.NewRequest(http.MethodGet, "/countryinfo/v1/info/nor", nil)
	req.SetPathValue("country_code", "nor")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
