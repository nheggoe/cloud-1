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
	PortRequired                = envRequiredErr(Port)
	CountryApiEndpointRequired  = envRequiredErr(CountriesEndpoint)
	CurrencyApiEndpointRequired = envRequiredErr(CurrencyEndpoint)
)

type Config struct {
	ServerSetting
	ApiEndpoint
}

type ServerSetting struct {
	Port string
}

type ApiEndpoint struct {
	Countries string
	Currency  string
}

func Load() (*Config, error) {
	config := &Config{
		ServerSetting{Port.GetOrDefault("8080")},
		ApiEndpoint{
			Countries: CountriesEndpoint.Get(),
			Currency:  CurrencyEndpoint.Get(),
		},
	}
	return config, validateConfig(config)
}

func validateConfig(cfg *Config) error {
	if cfg.Port == "" {
		return PortRequired
	}
	if cfg.Countries == "" {
		return CountryApiEndpointRequired
	}
	if cfg.Currency == "" {
		return CurrencyApiEndpointRequired
	}
	return nil
}

func envRequiredErr(env EnvVar) error {
	return fmt.Errorf("%s is required", env)
}
