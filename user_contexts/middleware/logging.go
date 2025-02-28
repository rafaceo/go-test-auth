package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
	"github.com/rafaceo/go-test-auth/user_contexts/service"
)

type loggingService struct {
	logger log.Logger
	next   service.UserContextService
}

func NewLoggingMiddleware(logger log.Logger, s service.UserContextService) service.UserContextService {
	return &loggingService{logger, s}
}

func (l *loggingService) AddUserContext(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, global bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "AddUserContext",
			"user_id", userID,
			"merchant_id", merchantID,
			"global", global,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.AddUserContext(ctx, userID, merchantID, global)
}

func (l *loggingService) EditUserContext(ctx context.Context, userID uuid.UUID, global bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "UpdateUserContext",
			"user_id", userID,
			"global", global,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.EditUserContext(ctx, userID, global)
}

func (l *loggingService) GetUserContexts(ctx context.Context, userID uuid.UUID) (contexts []domain.UserContext, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetUserContexts",
			"user_id", userID,
			"count", len(contexts),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.GetUserContexts(ctx, userID)
}

func (l *loggingService) DeleteUserContext(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "DeleteUserContext",
			"user_id", userID,
			"merchant_id", merchantID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.DeleteUserContext(ctx, userID, merchantID)
}

func (l *loggingService) DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "DeleteAllUserContexts",
			"user_id", userID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.DeleteAllUserContexts(ctx, userID)
}
