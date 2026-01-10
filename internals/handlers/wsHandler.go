package handlers

import (
	"log"
	"net/http"

	middleware "github.com/Ollefm/chat-backend-go/internals/middlewares"
	"github.com/Ollefm/chat-backend-go/internals/services"
	"github.com/Ollefm/chat-backend-go/internals/ws"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	auth        *services.AuthService
	chatService *services.ChatService
	hub         *ws.Hub
	upgrader    websocket.Upgrader
}

func NewWsHandler(auth *services.AuthService, chatSvc *services.ChatService, hub *ws.Hub) *WsHandler {
	return &WsHandler{
		auth:        auth,
		chatService: chatSvc,
		hub:         hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: Implement proper origin checking
			},
		},
	}
}

func (h *WsHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chatID", http.StatusBadRequest)
		return
	}

	// Upgrade connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	log.Printf("User %s connected to chat %s", userID, chatID)

	// Create and register client
	client := ws.NewClient(conn, userID, chatID)
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.ReadMessage(h.hub)
	go client.WriteMessage()

}
