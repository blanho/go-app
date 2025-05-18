// internal/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yourusername/azure-go-app/internal/models"
	"github.com/yourusername/azure-go-app/internal/service"
)

func RegisterRoutes(router *mux.Router, svc *service.Service) {
	router.HandleFunc("/users", getUsersHandler(svc)).Methods(http.MethodGet)
	router.HandleFunc("/users", createUserHandler(svc)).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", getUserHandler(svc)).Methods(http.MethodGet)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":500,"message":"Internal Server Error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, status int, message string, err error) {
	errorResponse := models.ErrorResponse{
		Status:  status,
		Message: message,
	}

	if err != nil {
		errorResponse.Error = err.Error()
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var details []models.ValidationError
		for _, err := range validationErrors {
			details = append(details, models.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		errorResponse.Details = details
	}

	respondWithJSON(w, status, errorResponse)
}