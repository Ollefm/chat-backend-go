package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	middleware "github.com/Ollefm/chat-backend-go/internals/middlewares"
	"github.com/Ollefm/chat-backend-go/internals/services"
)

type ChatHandler struct {
	chatService *services.ChatService
	userService *services.UserService
}

func NewChatHandler(chatService *services.ChatService, userService *services.UserService) *ChatHandler {
	return &ChatHandler{chatService, userService}
}

func (h *ChatHandler) InitiateChat(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse JSON body
	var body struct {
		TargetUserID string `json:"target_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.TargetUserID == "" {
		http.Error(w, "missing target_user_id in request body", http.StatusBadRequest)
		return
	}

	// Initiate chat via service
	chat, err := h.chatService.InitiateChat(r.Context(), userID, body.TargetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return chat JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	chatID := r.URL.Query().Get("chat_id")
	history, err := h.chatService.GetHistory(r.Context(), userID, chatID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(history)
}

func (h *ChatHandler) GetChatData(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	chatID := r.URL.Query().Get("chat_id")

	chatInfo, err := h.chatService.GetChatData(r.Context(), chatID)
	if err != nil {
		log.Printf("GetChatData error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(chatInfo)
}
