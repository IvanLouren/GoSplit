package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/IvanLouren/GoSplit/pkg/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type UpdateMeRequest struct {
	Name string `json:"name"`
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetMe(userID)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var req UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name must not be empty", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateMe(userID, req.Name)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}
