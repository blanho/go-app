package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/azure-go-app/internal/config"
	"github.com/yourusername/azure-go-app/internal/handlers"
	"github.com/yourusername/azure-go-app/internal/middleware"
	"github.com/yourusername/azure-go-app/internal/repository"
	"github.com/yourusername/azure-go-app/internal/service"
	"github.com/yourusername/azure-go-app/internal/telemetry"
)

func main() {
	ctx := context.Background()
	
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	var dbPassword string
	if cfg.Environment != "local" {
		kv, err := config.NewKeyVaultClient()
		if err != nil {
			log.Fatalf("Failed to initialize Key Vault client: %v", err)
		}
		
		dbPassword, err = kv.GetSecret(ctx, "DB-PASSWORD")
		if err != nil {
			log.Fatalf("Failed to get database password: %v", err)
		}
		
		cfg.DatabaseURL = cfg.DatabaseURL + ";Password=" + dbPassword
	}
	
	tel := telemetry.NewTelemetry(cfg)
	defer tel.Flush()
	
	repo, err := repository.NewRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()
	
	svc := service.NewService(repo, tel)
	
	router := mux.NewRouter()
	
	router.Use(middleware.TelemetryMiddleware(tel))
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware(tel))
	
	router.HandleFunc("/health", handlers.HealthCheckHandler(repo)).Methods(http.MethodGet)
	router.HandleFunc("/ready", handlers.ReadinessCheckHandler(repo)).Methods(http.MethodGet)
	router.HandleFunc("/metrics", handlers.MetricsHandler()).Methods(http.MethodGet)
	
	apiRouter := router.PathPrefix("/api").Subrouter()
	handlers.RegisterRoutes(apiRouter, svc)
	
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	go func() {
		log.Printf("Starting server on port %s in %s environment", cfg.Port, cfg.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	
	tel.TrackEvent("ServiceShutdown", map[string]string{
		"reason": "graceful",
	}, nil)
	
	if err := server.Shutdown(ctx); err != nil {
		tel.TrackException(err, map[string]string{
			"component": "server",
			"operation": "shutdown",
		})
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server gracefully stopped")
}