package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-starter/internal/config"
	"go-starter/internal/handlers"
	"go-starter/internal/logger"
	"go-starter/internal/middleware"
	"go-starter/internal/repositories"
	"go-starter/internal/services"
	"go-starter/pkg/database"

	_ "go-starter/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Go Starter API
// @version 1.0
// @description RESTful API with JWT authentication, rate limiting, and PostgreSQL
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logger.Level, cfg.IsProduction()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting application",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.Server.Port),
	)

	// Initialize database
	db, err := database.New(database.Config{
		DSN:             cfg.GetDSN(),
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}, logger.Get())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler(db)

	// Create router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware(cfg.IsProduction()))
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit.RPS, cfg.RateLimit.Burst))

	// Health check routes (no auth required)
	router.HandleFunc("/healthz", healthHandler.Healthz).Methods("GET")
	router.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	// Auth routes (no auth required)
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Swagger documentation (only in development)
	if !cfg.IsProduction() {
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		logger.Info("swagger documentation enabled at /swagger/index.html")
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}
