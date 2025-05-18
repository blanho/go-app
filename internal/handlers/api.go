// internal/handlers/api.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/azure-go-app/internal/models"
	"github.com/yourusername/azure-go-app/internal/service"
)

func getUserHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		user, err := svc.GetUserByID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
			return
		}

		if user == nil {
			respondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}

		response := models.NewResponse(http.StatusOK, "User retrieved successfully", user)
		respondWithJSON(w, http.StatusOK, response)
	}
}

func createUserHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.UserInput
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&input); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
			return
		}
		defer r.Body.Close()

		user, err := svc.CreateUser(r.Context(), input)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		response := models.NewResponse(http.StatusCreated, "User created successfully", user)
		respondWithJSON(w, http.StatusCreated, response)
	}
}

func getUsersHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("pageSize")

		page := 1
		if pageStr != "" {
			pageNum, err := strconv.Atoi(pageStr)
			if err == nil && pageNum > 0 {
				page = pageNum
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			size, err := strconv.Atoi(pageSizeStr)
			if err == nil && size > 0 && size <= 100 {
				pageSize = size
			}
		}

		users, err := svc.ListUsers(r.Context(), page, pageSize)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve users", err)
			return
		}

		response := models.NewResponse(http.StatusOK, "Users retrieved successfully", users)
		respondWithJSON(w, http.StatusOK, response)
	}
}