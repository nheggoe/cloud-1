package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	ApiEndpoint
}

type ApiEndpoint struct {
	Countries string
	Currency  string
}

func Load() (*Config, error) {
	config := &Config{
		ApiEndpoint{
			Countries: os.Getenv("COUNTRIES_ENDPOINT"),
			Currency:  os.Getenv("CURRENCY_ENDPOINT"),
		},
	}
	return config, validateConfig(config)
}

func validateConfig(cfg *Config) error {
	var out strings.Builder
	if cfg.Countries == "" {
		out.WriteString("COUNTRIES_ENDPOINT is required")
	}
	if cfg.Currency == "" {
		if out.Len() != 0 {
			out.WriteByte('\n')
		}
		out.WriteString("CURRENCY_ENDPOINT is required")
	}
	if out.Len() != 0 {
		return errors.New(out.String())
	}
	return nil
}
