package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
	"github.com/rafaceo/go-test-auth/user_contexts/repository"
)

type UserContextService interface {
	AddUserContext(ctx context.Context, userID uuid.UUID, merchantID string, global bool) error
	EditUserContext(ctx context.Context, userID uuid.UUID, global bool) error
	GetUserContexts(ctx context.Context, userID uuid.UUID) (domain.UserContext, error)
	DeleteUserContext(ctx context.Context, userID uuid.UUID, merchantID string) error
	DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) error
}

type userContextService struct {
	repo repository.UserContextRepository
}

func NewUserContextService(repo repository.UserContextRepository) UserContextService {
	return &userContextService{repo: repo}
}

func (s *userContextService) AddUserContext(ctx context.Context, userID uuid.UUID, merchantID string, global bool) error {
	if global {
		return s.repo.EditUserContext(ctx, userID, global)
	}
	return s.repo.AddUserContext(ctx, domain.UserContext{
		UserID:     userID,
		MerchantID: merchantID, // Конвертируем в строку
		Global:     global,     // Передаём глобальный флаг
	})
}

func (s *userContextService) EditUserContext(ctx context.Context, userID uuid.UUID, global bool) error {
	return s.repo.EditUserContext(ctx, userID, global)
}

func (s *userContextService) GetUserContexts(ctx context.Context, userID uuid.UUID) (domain.UserContext, error) {
	return s.repo.GetUserContexts(ctx, userID)
}

func (s *userContextService) DeleteUserContext(ctx context.Context, userID uuid.UUID, merchantID string) error {
	return s.repo.DeleteUserContext(ctx, userID, merchantID)
}

func (s *userContextService) DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteAllUserContexts(ctx, userID)
}
