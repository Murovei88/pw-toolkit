package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/murovei88/pw-toolkit/internal/config"
	"github.com/murovei88/pw-toolkit/internal/handler"
	"github.com/murovei88/pw-toolkit/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting PW Toolkit API",
		"version", cfg.Version,
		"port", cfg.Port,
	)

	// Setup router and handlers
	router := http.NewServeMux()
	
	// Health check endpoints
	router.HandleFunc("GET /api/v1/status", handler.StatusHandler(cfg))
	router.HandleFunc("GET /api/v1/health", handler.HealthHandler())

	// Apply middleware
	handler := middleware.Chain(
		router,
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.CORS(),
	)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}
