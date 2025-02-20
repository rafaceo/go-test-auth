package middleware

import (
	"context"
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

func (l *loggingService) CreateUser(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "CreateUser",
				"took", time.Since(begin),
				"context", ctx,
				"name", name,
				"err", err,
			)
		}
	}(time.Now())

	return l.next.CreateUser(ctx, name)
}

func (l *loggingService) EditUser(ctx context.Context, id uint, name, context string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "EditUser",
				"took", time.Since(begin),
				"context", context,
				"id", id,
				"name", name,
				"err", err,
			)
		}
	}(time.Now())

	return l.next.EditUser(ctx, id, name, context)
}
func (l *loggingService) GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (l *loggingService) EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (l *loggingService) RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (l *loggingService) GetUser(ctx context.Context, id uint) (name string, context string, rights map[string][]string, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GetUser",
				"id", id,
				"name", name,
				"context", context,
				"rights", rights,
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	name, context, rights, err = l.next.GetUser(ctx, id)
	return
}
