package services

import (
	"context"
	"errors"
	"time"

	"github.com/Ollefm/chat-backend-go/internals/models"
	"github.com/Ollefm/chat-backend-go/internals/repositories"
)

type ChatService struct {
	chatRepo *repositories.ChatRepo
	msgRepo  *repositories.MessageRepo
}

func NewChatService(
	cRepo *repositories.ChatRepo,
	mRepo *repositories.MessageRepo,
) *ChatService {
	return &ChatService{
		chatRepo: cRepo,
		msgRepo:  mRepo,
	}
}

func (s *ChatService) InitiateChat(ctx context.Context, userID string, other string) (*models.Chat, error) {

	return s.chatRepo.GetOrCreateChat(ctx, userID, other)
}

func (s *ChatService) GetHistory(ctx context.Context, userID string, chatID string) ([]models.Message, error) {

	ok, err := s.chatRepo.UserInChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("forbidden")
	}

	return s.msgRepo.GetMessages(ctx, chatID)
}

func (s *ChatService) GetChatData(ctx context.Context, chatID string) (*models.Chat, error) {
	return s.chatRepo.GetChatData(ctx, chatID)
}

func (s *ChatService) SaveMessage(ctx context.Context, chatID string, senderID string, content string) (*models.Message, []string, error) {

	ok, err := s.chatRepo.UserInChat(ctx, senderID, chatID)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, errors.New("forbidden")
	}

	msg := &models.Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.msgRepo.SaveMessage(ctx, msg); err != nil {

		return nil, nil, err
	}

	participants, err := s.chatRepo.GetParticipants(ctx, chatID)
	if err != nil {
		return nil, nil, err
	}

	return msg, participants, nil
}
