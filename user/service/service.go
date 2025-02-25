package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/user/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"strings"
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
	CreateUser(ctx context.Context, phone string, passwordHash string) error
	EditUser(ctx context.Context, id uuid.UUID, phone, password string) error
	GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	GetUser(ctx context.Context, id uuid.UUID) (string, string, string, string, error)
	GetUserRights(ctx context.Context, id uuid.UUID) (map[string][]string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, phone string, passwordHash string) error {
	re := regexp.MustCompile(`^\+?\d{10,}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number: must contain only digits, optionally start with '+', and have at least 10 digits")
	}
	strings.ReplaceAll(phone, " ", "")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, phone, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) EditUser(ctx context.Context, id uuid.UUID, phone, password string) error {

	if phone == "" || password == "" {
		return errors.New("phone and password cannot be empty")
	}

	re := regexp.MustCompile(`^\+?\d{10,}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number: must contain only digits, optionally start with '+', and have at least 10 digits")
	}
	strings.ReplaceAll(phone, " ", "")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return errors.New("failed to hash password")
	}

	err = s.repo.EditUser(ctx, id, phone, string(hashedPassword))
	if err != nil {
		log.Printf("Error editing user: %v", err)
		return err
	}

	return nil
}

func (s *userService) GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
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

func (s *userService) EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
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

func (s *userService) RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
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

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (string, string, string, string, error) {
	phone, passwordHash, createdAt, updatedAt, err := s.repo.GetUser(ctx, id)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return "", "", "", "", err
	}
	return phone, passwordHash, createdAt, updatedAt, nil
}

func (s *userService) GetUserRights(ctx context.Context, id uuid.UUID) (map[string][]string, error) {
	rights, err := s.repo.GetUserRights(ctx, id)
	if err != nil {
		log.Printf("Error retrieving user rights: %v", err)
		return nil, err
	}
	return rights, nil
}
