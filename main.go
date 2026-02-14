package main

import (
	"cloud1/config"
	"cloud1/router"
	"context"
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

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading environment variable", "error", err)
	}

	mux := router.NewRouter(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", cfg.Port),
		Handler: router.LoggingMiddleware(mux),
	}
	go func() {
		slog.Info("Server is live", "addr", server.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	shutdown(server)
}

func shutdown(server *http.Server) {
	slog.Info("Gracefully shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}
