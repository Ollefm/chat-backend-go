package services

import (
	"context"
	"errors"
	"time"

	"github.com/Ollefm/chat-backend-go/internals/repositories"
	"github.com/Ollefm/chat-backend-go/pkg/hashing"
	"github.com/Ollefm/chat-backend-go/pkg/helpers"
)

type AuthService struct {
	userRepo    *repositories.UserRepo
	sessionRepo *repositories.SessionRepo
}

func NewAuthService(uRepo *repositories.UserRepo, tRepo *repositories.SessionRepo) *AuthService {
	return &AuthService{uRepo, tRepo}
}

func generateSessionID() (string, error) {
	return helpers.GenerateSessionToken(64) // 256 bits
}

func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if len(username) > 30 {
		return "", errors.New("username cannot exceed the maximum length of 30 characters")
	}
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}

	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return "", err
	}

	return s.userRepo.CreateUser(ctx, username, string(hashedPassword))
}

func (s *AuthService) Username(ctx context.Context, id string) (string, error) {
	username, err := s.userRepo.GetUsernameByUid(ctx, id)
	if err != nil {
		return "", errors.New("no user found")
	}
	return username, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	userID, hash, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	valid, err := hashing.CheckPasswordHash(password, hash)
	if err != nil || !valid {
		return "", "", errors.New("invalid credentials")
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return "", "", err
	}

	sessionHash := hashing.HashToken(sessionID)

	err = s.sessionRepo.CreateSession(
		ctx,
		sessionHash,
		userID,
		time.Now().Add(30*24*time.Hour), // 30 days
	)
	if err != nil {
		return "", "", err
	}

	return sessionID, userID, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	hash := hashing.HashToken(sessionID)
	return s.sessionRepo.DeleteSession(ctx, hash)
}

func (s *AuthService) GetUserBySession(ctx context.Context, sessionID string) (string, string, error) {
	if sessionID == "" {
		return "", "", errors.New("missing session")
	}

	hash := hashing.HashToken(sessionID)

	userID, expiresAt, err := s.sessionRepo.GetSession(ctx, hash)
	if err != nil {
		return "", "", errors.New("invalid session")
	}

	username, err := s.userRepo.GetUsernameByUid(ctx, userID)
	if err != nil {
		return "", "", errors.New("no username found")
	}

	if time.Now().After(expiresAt) {
		_ = s.sessionRepo.DeleteSession(ctx, hash)
		return "", "", errors.New("session expired")
	}

	return userID, username, nil
}
