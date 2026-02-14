package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	Port              = "PORT"
	CountriesEndpoint = "COUNTRIES_ENDPOINT"
	CurrencyEndpoint  = "CURRENCY_ENDPOINT"
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
	Port int
}

type ApiEndpoint struct {
	Countries string
	Currency  string
}

func Load() (*Config, error) {

	config := &Config{
		ServerSetting{GetenvOrDefault(Port, 8080, func(s string) (int, error) { return strconv.Atoi(s) })},
		ApiEndpoint{
			Countries: os.Getenv(CountriesEndpoint),
			Currency:  os.Getenv(CurrencyEndpoint),
		},
	}
	return config, validateConfig(config)
}

func validateConfig(cfg *Config) error {
	if cfg.Port == 0 {
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

func GetenvOrDefault[T any](name string, defaultValue T, f func(string) (T, error)) T {
	if val, err := f(os.Getenv(name)); err == nil {
		return val
	}
	return defaultValue
}

func envRequiredErr(env string) error {
	return fmt.Errorf("%s is required", env)
}
