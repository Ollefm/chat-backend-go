package models

import (
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// ----------- Chat STRUCTS -----------
type Chat struct {
	ID            string    `json:"id"`
	User1ID       string    `json:"user1_id"`
	User1Username string    `json:"username1"`
	User2ID       string    `json:"user2_id"`
	User2Username string    `json:"username2"`
	Created       time.Time `json:"created_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
