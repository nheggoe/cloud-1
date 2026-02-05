package main

import (
	"cloud1/config"
	"cloud1/router"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
		Addr:    "localhost:8080",
		Handler: router.LoggingMiddleware(mux),
	}
	defer slog.Info("Gracefully shutting down")
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}

func setupLogger() {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level.Set(slog.LevelInfo)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	var handler slog.Handler
	switch strings.ToLower(os.Getenv("LOG_FORMAT")) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
