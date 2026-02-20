package config

import (
	"fmt"
	"os"
)

type EnvVar string

func (e EnvVar) Get() string {
	return os.Getenv(string(e))
}

func (e EnvVar) GetOrDefault(defaultValue string) string {
	if res := e.Get(); res != "" {
		return res
	}
	return defaultValue
}

const (
	Port              EnvVar = "PORT"
	CountriesEndpoint EnvVar = "COUNTRIES_ENDPOINT"
	CurrencyEndpoint  EnvVar = "CURRENCY_ENDPOINT"
)

var (
	CountryAPIEndpointRequired  = envRequiredErr(CountriesEndpoint)
	CurrencyAPIEndpointRequired = envRequiredErr(CurrencyEndpoint)
)

type Config struct {
	ServerSetting
	APIEndpoint
}

type ServerSetting struct {
	Port string
}

type APIEndpoint struct {
	CountriesEndpoint string
	CurrencyEndpoint  string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerSetting{Port.GetOrDefault("8080")},
		APIEndpoint{
			CountriesEndpoint: CountriesEndpoint.Get(),
			CurrencyEndpoint:  CurrencyEndpoint.Get(),
		},
	}
	return cfg, validateConfig(cfg)
}

func validateConfig(cfg *Config) error {
	if cfg.CountriesEndpoint == "" {
		return CountryAPIEndpointRequired
	}
	if cfg.CurrencyEndpoint == "" {
		return CurrencyAPIEndpointRequired
	}
	return nil
}

func envRequiredErr(env EnvVar) error {
	return fmt.Errorf("%s is required", env)
}
