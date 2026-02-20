package info

import (
	"countryinfo/internal/restclient"
	"countryinfo/internal/util"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type service struct {
	countries *restclient.CountriesClient
}

type Response struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Area       int               `json:"area"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
}

func NewResponse(c restclient.Country) Response {
	return Response{
		Name:       c.Name.Common,
		Continents: c.Continents,
		Population: c.Population,
		Area:       int(c.Area),
		Languages:  c.Languages,
		Borders:    c.Borders,
		Flag:       c.Flags.Png,
		Capital:    c.Capital[0],
	}
}

func Handler(countries *restclient.CountriesClient) http.HandlerFunc {
	s := &service{
		countries: countries,
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

	countries, err := s.countries.GetByAlpha(r.Context(), countryCode)
	if err != nil {
		slog.ErrorContext(r.Context(), "upstream countries request failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	resJson, err := json.Marshal(NewResponse(countries[0]))
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resJson)

	slog.InfoContext(
		r.Context(),
		"country info request completed",
		"country_code", countryCode,
	)
}
