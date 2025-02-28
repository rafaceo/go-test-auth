package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/rafaceo/go-test-auth/rights/domain"
	"github.com/rafaceo/go-test-auth/rights/service"
	"time"
)

type loggingService struct {
	logger log.Logger
	next   service.RightsService
}

// AddRights
func (l *loggingService) AddRights(ctx context.Context, module string, action []string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "AddRights",
			"module", module,
			"action", action,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.AddRights(ctx, module, action)
}

// EditRight
func (l *loggingService) EditRight(ctx context.Context, id, module string, action []string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "EditRight",
			"id", id,
			"module", module,
			"action", action,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.EditRight(ctx, id, module, action)
}

// GetAllRights
func (l *loggingService) GetAllRights(ctx context.Context) (rights []domain.Right, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetAllRights",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.GetAllRights(ctx)
}

// GetRightByName
func (l *loggingService) GetRightByName(ctx context.Context, module string) (right *domain.Right, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetRightByName",
			"module", module,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.GetRightByName(ctx, module)
}

// GetRightById
func (l *loggingService) GetRightById(ctx context.Context, id string) (right *domain.Right, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetRightById",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.GetRightById(ctx, id)
}

// DeleteRight
func (l *loggingService) DeleteRight(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "DeleteRight",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return l.next.DeleteRight(ctx, id)
}

func NewLoggingMiddleware(logger log.Logger, next service.RightsService) *loggingService {
	return &loggingService{
		logger: logger,
		next:   next,
	}
}
