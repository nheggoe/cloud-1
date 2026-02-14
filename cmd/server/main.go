package main

import (
	"countryinfo/internal/config"
	"countryinfo/internal/server"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}

	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
