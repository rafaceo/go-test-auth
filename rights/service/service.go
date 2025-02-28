package service

import (
	"context"
	"fmt"
	"github.com/rafaceo/go-test-auth/rights/domain"
	"github.com/rafaceo/go-test-auth/rights/repository"
)

type rightsService struct {
	repo repository.RightsRepository
}

type RightsService interface {
	AddRights(ctx context.Context, module string, action []string) error
	EditRight(ctx context.Context, id string, module string, action []string) error
	GetAllRights(ctx context.Context) ([]domain.Right, error)
	GetRightByName(ctx context.Context, module string) (*domain.Right, error)
	GetRightById(ctx context.Context, id string) (*domain.Right, error)
	DeleteRight(ctx context.Context, id string) error
}

func NewRightsService(repo repository.RightsRepository) RightsService {
	return &rightsService{repo: repo}
}

func (s *rightsService) AddRights(ctx context.Context, module string, action []string) error {
	return s.repo.AddRights(ctx, module, action)
}

func (s *rightsService) EditRight(ctx context.Context, id string, module string, action []string) error {
	return s.repo.EditRight(ctx, id, module, action)
}

func (s *rightsService) GetAllRights(ctx context.Context) ([]domain.Right, error) {
	return s.repo.GetAllRights(ctx)
}

func (s *rightsService) GetRightByName(ctx context.Context, module string) (*domain.Right, error) {
	return s.repo.GetRightByName(ctx, module)
}

func (s *rightsService) GetRightById(ctx context.Context, id string) (*domain.Right, error) {
	fmt.Println("Searching for Right ID:", id)
	return s.repo.GetRightById(ctx, id)
}

func (s *rightsService) DeleteRight(ctx context.Context, id string) error {
	return s.repo.DeleteRight(ctx, id)
}
