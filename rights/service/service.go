package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rafaceo/go-test-auth/rights/repository"
	"log"
	"regexp"
)

var allowedSections = map[string]bool{
	"ORDERS":        true,
	"PRODUCTS_BASE": true,
	"PRODUCTS":      true,
	"PRICELIST":     true,
	"PRODUCTS_ADD":  true,
	"CATALOG":       true,
	"CHAR":          true,
	"HISTORY":       true,
}

var allowedPermissions = map[string]bool{
	"READ":   true,
	"CREATE": true,
	"UPDATE": true,
	"DELETE": true,
}

type UserService interface {
	CreateUser(ctx context.Context, name string) error
	EditUser(ctx context.Context, id uint, name, context string) error
	GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) error
	EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) error
	RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) error
	GetUser(ctx context.Context, id uint) (string, string, map[string][]string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, name string) error {
	fmt.Println()
	if name == "" {
		return errors.New("name cannot be fuck")
	}

	err := s.repo.CreateUser(ctx, name)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

func (s *userService) EditUser(ctx context.Context, id uint, name, context string) error {
	if name == "" || context == "" {
		return errors.New("name and context cannot be empty")
	}

	merchantRegex := regexp.MustCompile(`^MERCHANT_\d+$`)

	if context != "MERCHANT_ALL" && !merchantRegex.MatchString(context) {
		return errors.New("invalid context: must be MERCHANT_ALL or MERCHANT_{merchant_id}")
	}

	err := s.repo.EditUser(ctx, id, name, context)
	if err != nil {
		log.Printf("Error editing user: %v", err)
		return err
	}

	return nil
}

func (s *userService) GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) error {
	if len(rights) == 0 {
		return errors.New("rights cannot be empty")
	}

	// Проверяем права
	for section, perms := range rights {
		// Проверяем, есть ли такая секция в списке разрешённых
		if !allowedSections[section] {
			return errors.New("invalid section: " + section)
		}

		// Проверяем, что все права в списке разрешены
		for _, perm := range perms {
			if !allowedPermissions[perm] {
				return errors.New("invalid permission: " + perm)
			}
		}
	}

	err := s.repo.GrantRightsToUser(ctx, id, rights)
	if err != nil {
		log.Printf("Error granting rights to user: %v", err)
		return err
	}

	return nil
}

func (s *userService) EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) error {
	if len(rights) == 0 {
		return errors.New("rights cannot be empty")
	}

	for section, perms := range rights {
		if !allowedSections[section] {
			return errors.New("invalid section: " + section)
		}

		for _, perm := range perms {
			if !allowedPermissions[perm] {
				return errors.New("invalid permission: " + perm)
			}
		}
	}

	err := s.repo.EditRightsToUser(ctx, id, rights)
	if err != nil {
		log.Printf("Error granting rights to user: %v", err)
		return err
	}

	return nil
}

func (s *userService) RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) error {
	if len(rights) == 0 {
		return errors.New("rights cannot be empty")
	}

	for section, perms := range rights {
		if !allowedSections[section] {
			return errors.New("invalid section: " + section)
		}

		for _, perm := range perms {
			if !allowedPermissions[perm] {
				return errors.New("invalid permission: " + perm)
			}
		}
	}

	err := s.repo.RevokeRightsFromUser(ctx, id, rights)
	if err != nil {
		log.Printf("Error granting rights to user: %v", err)
		return err
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, id uint) (string, string, map[string][]string, error) {
	name, context, rights, err := s.repo.GetUser(ctx, id)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return "", "", nil, err
	}
	return name, context, rights, nil
}
