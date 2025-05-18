// internal/handlers/health.go
package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yourusername/azure-go-app/internal/models"
	"github.com/yourusername/azure-go-app/internal/repository"
)

func HealthCheckHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.CheckHealth(r.Context())
		
		if err != nil {
			respondWithError(w, http.StatusServiceUnavailable, "Service unavailable", err)
			return
		}
		
		response := models.NewResponse(http.StatusOK, "Service healthy", nil)
		respondWithJSON(w, http.StatusOK, response)
	}
}

func ReadinessCheckHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := repo.HealthStatus(r.Context())
		
		if !status.Database {
			respondWithError(w, http.StatusServiceUnavailable, "Service not ready", nil)
			return
		}
		
		response := models.NewResponse(http.StatusOK, "Service ready", status)
		respondWithJSON(w, http.StatusOK, response)
	}
}


func MetricsHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}