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
	"github.com/murovei88/pw-toolkit/internal/database"
	"github.com/murovei88/pw-toolkit/internal/handler"
	"github.com/murovei88/pw-toolkit/internal/middleware"
	"github.com/murovei88/pw-toolkit/internal/repository/mysql"
	"github.com/murovei88/pw-toolkit/internal/service"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting PW Toolkit API",
		"version", cfg.Version,
		"port", cfg.Port,
	)

	// Initialize database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connected", "host", cfg.DBHost, "port", cfg.DBPort)

	// Initialize Redis
	rdb, err := database.NewRedisClient(cfg)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	logger.Info("Redis connected", "host", cfg.RedisHost, "port", cfg.RedisPort)

	// Initialize repositories
	buildRepo := mysql.NewBuildRepository(db)
	classRepo := mysql.NewClassRepository(db)
	itemRepo := mysql.NewItemRepository(db)
	gemRepo := mysql.NewGemRepository(db)

	// Initialize services
	buildService := service.NewBuildService(buildRepo, classRepo, itemRepo, gemRepo)

	// Initialize handlers
	statusHandler := handler.StatusHandler(cfg)
	healthHandler := handler.HealthHandler(db, rdb)
	buildHandler := handler.NewBuildHandler(buildService, logger)
	calcHandler := handler.NewCalculationHandler(buildService, logger)

	// Setup router
	router := http.NewServeMux()
	
	// Health/status endpoints
	router.HandleFunc("GET /api/v1/status", statusHandler)
	router.HandleFunc("GET /api/v1/health", healthHandler)

	// Build endpoints
	router.HandleFunc("POST /api/v1/builds", buildHandler.CreateBuild)
	router.HandleFunc("GET /api/v1/builds/{id}", buildHandler.GetBuild)

	// Calculation endpoint
	router.HandleFunc("POST /api/v1/calculate", calcHandler.CalculatePreview)

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
