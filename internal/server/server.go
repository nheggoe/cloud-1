package server

import (
	"context"
	"countryinfo/internal/config"
	"countryinfo/internal/middleware"
	"countryinfo/internal/router"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server represents the HTTP server with its configuration.
type Server struct {
	http *http.Server
}

// New creates and configures a new HTTP server.
func New(cfg *config.Config) *Server {
	mux := router.New(cfg)

	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: middleware.Logging(mux),
		},
	}
}

// Run starts the server and blocks until it receives a shutdown signal.
// It handles graceful shutdown automatically.
func (s *Server) Run() error {
	// Start server in a goroutine
	go func() {
		slog.Info("Server started", "addr", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	return s.shutdown()
}

// shutdown gracefully shuts down the server with a timeout.
func (s *Server) shutdown() error {
	slog.Info("Gracefully shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}
