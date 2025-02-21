package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/rights/service"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	next   service.UserService
}

func NewLoggingMiddleware(logger log.Logger, next service.UserService) service.UserService {
	return &loggingService{logger: logger, next: next}
}

func (l *loggingService) CreateUser(ctx context.Context, phone string, passwordHash string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "CreateUser",
				"took", time.Since(begin),
				"context", ctx,
				"phone", phone,
				"passwordHash", passwordHash,
				"err", err,
			)
		}
	}(time.Now())

	return l.next.CreateUser(ctx, phone, passwordHash)
}

func (l *loggingService) EditUser(ctx context.Context, id uuid.UUID, phone, password string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "EditUser",
				"took", time.Since(begin),
				"id", id,
				"phone", phone,
				"password", password,
				"err", err,
			)
		}
	}(time.Now())

	return l.next.EditUser(ctx, id, phone, password)
}
func (l *loggingService) GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GrantRightsToUser",
				"took", time.Since(begin),
				"context", ctx,
				"id", id,
				"rights", rights,
				"err", err,
			)
		}

	}(time.Now())
	return l.next.GrantRightsToUser(ctx, id, rights)
}

func (l *loggingService) EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "EditRightsToUser",
				"took", time.Since(begin),
				"context", ctx,
				"id", id,
				"rights", rights,
				"err", err,
			)
		}
	}(time.Now())
	return l.next.EditRightsToUser(ctx, id, rights)
}

func (l *loggingService) RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "RevokeRightsFromUser",
				"took", time.Since(begin),
				"context", ctx,
				"id", id,
				"rights", rights,
				"err", err,
			)
		}
	}(time.Now())
	return l.next.RevokeRightsFromUser(ctx, id, rights)
}

func (l *loggingService) GetUser(ctx context.Context, id uuid.UUID) (phone string, password string, created_at string, updated_at string, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GetUser",
				"id", id,
				"phone", phone,
				"password", password,
				"createdAt", created_at,
				"updatedAt", updated_at,
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	phone, password, created_at, updated_at, err = l.next.GetUser(ctx, id)
	return
}

func (l *loggingService) GetUserRights(ctx context.Context, id uuid.UUID) (rights map[string][]string, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GetUserRights",
				"id", id,
				"took", time.Since(begin),
				"rights", rights,
				"err", err,
			)
		}
	}(time.Now())
	rights, err = l.next.GetUserRights(ctx, id)
	return
}
