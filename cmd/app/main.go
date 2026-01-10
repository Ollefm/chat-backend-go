package main

import (
	"log"
	"net/http"

	"github.com/Ollefm/chat-backend-go/internals/databases"
	"github.com/Ollefm/chat-backend-go/internals/env"
	"github.com/Ollefm/chat-backend-go/internals/handlers"
	middleware "github.com/Ollefm/chat-backend-go/internals/middlewares"
	"github.com/Ollefm/chat-backend-go/internals/repositories"
	"github.com/Ollefm/chat-backend-go/internals/services"
	"github.com/Ollefm/chat-backend-go/internals/ws"
)

func main() {

	// Load .env
	if err := env.Load(); err != nil {
		log.Println("Warning: .env file not found. Using system ENV variables.")
	}

	// Connect database
	db, err := databases.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repos
	userRepo := repositories.NewUserRepo(db)
	sessionRepo := repositories.NewSessionRepo(db)
	messageRepo := repositories.NewMessageRepo(db)
	chatRepo := repositories.NewChatRepo(db)

	// Create and start the hub
	hub := ws.NewHub()
	go hub.Run()

	// Services
	authService := services.NewAuthService(userRepo, sessionRepo)
	chatService := services.NewChatService(chatRepo, messageRepo)
	userService := services.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	chatHandler := handlers.NewChatHandler(chatService, userService)
	wsHandler := handlers.NewWsHandler(authService, chatService, hub)
	userHandler := handlers.NewUserHandler(userService)

	// PUBLIC ROUTES (no auth required)
	r := http.NewServeMux()
	r.HandleFunc("POST /auth/register/", authHandler.RegisterHandler)
	r.HandleFunc("POST /auth/login/", authHandler.LoginHandler)
	r.HandleFunc("POST /auth/logout/", authHandler.LogoutHandler)

	// PROTECTED ROUTES (auth middleware used)
	r.Handle("GET /api/ws", middleware.AuthMiddleware(authService)(http.HandlerFunc(wsHandler.ServeWS)))

	r.Handle("GET /api/users/search", middleware.AuthMiddleware(authService)(
		http.HandlerFunc(userHandler.SearchUsers),
	))
	r.Handle("GET /api/users/me/", middleware.AuthMiddleware(authService)(
		http.HandlerFunc(authHandler.MeHandler),
	))

	r.Handle("POST /api/chat/initiatechat", middleware.AuthMiddleware(authService)(
		http.HandlerFunc(chatHandler.InitiateChat),
	))
	r.Handle("GET /api/chat/history/", middleware.AuthMiddleware(authService)(
		http.HandlerFunc(chatHandler.GetHistory),
	))
	r.Handle("GET /api/chat/chatinfo/", middleware.AuthMiddleware(authService)(
		http.HandlerFunc(chatHandler.GetChatData),
	))

	// CORS
	var allowed = map[string]bool{
		"http://localhost:3000": true,
		//"https://produrl": true,
	}

	serverHandler := middleware.CORS(allowed)(r)
	// Start server
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", serverHandler); err != nil {
		log.Fatal(err)
	}
}
