package services

import (
	"context"

	"github.com/Ollefm/chat-backend-go/internals/models"
	"github.com/Ollefm/chat-backend-go/internals/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepo
}

func NewUserService(uRepo *repositories.UserRepo) *UserService {
	return &UserService{userRepo: uRepo}
}

func (s *UserService) SearchUsers(ctx context.Context, query string, excludeUserID string) ([]models.UserResult, error) {
	res, err := s.userRepo.SearchUsers(ctx, query, excludeUserID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
