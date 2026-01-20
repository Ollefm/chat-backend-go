# Chat Backend - Go (Work in Progress)

A real-time chat application backend built with Go, featuring WebSocket support, PostgreSQL database, and session-based authentication.

## Features

- **User Authentication**: Register, login, and logout with secure password hashing
- **Real-time Messaging**: WebSocket support for instant message delivery
- **Chat Management**: Initiate chats, retrieve chat history
- **User Search**: Search for users to initiate conversations
- **Session Management**: Secure session-based authentication
- **CORS Support**: Configured cross-origin for frontend integration

## Tech Stack

- **Language**: Go 1.25.4
- **Database**: PostgreSQL 16
- **WebSocket**: Gorilla WebSocket
- **Cryptography**: golang.org/x/crypto
- **Database Driver**: jackc/pgx v5
- **Containerization**: Docker & Docker Compose
