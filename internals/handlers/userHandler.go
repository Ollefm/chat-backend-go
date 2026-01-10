package handlers

import (
	"encoding/json"

	"net/http"

	middleware "github.com/Ollefm/chat-backend-go/internals/middlewares"
	"github.com/Ollefm/chat-backend-go/internals/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing 'q' query parameter", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	results, err := h.userService.SearchUsers(r.Context(), q, userID)
	if err != nil {
		http.Error(w, "failed to search users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
