package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/service"
	"time"
)

type loggingService struct {
	logger log.Logger
	next   service.RoleService
}

func (l loggingService) AddRole(ctx context.Context, Role domain.Role) (id uuid.UUID, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "AddRole",
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	return l.next.AddRole(ctx, Role)
}

func (l loggingService) UpdateRole(ctx context.Context, Role domain.Role) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "UpdateRole",
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	return l.next.UpdateRole(ctx, Role)
}

func (l loggingService) GetRoleById(ctx context.Context, id uuid.UUID) (Role *domain.Role, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GetRoleById",
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	return l.next.GetRoleByID(ctx, id)
}

func (l loggingService) GetAllRoles(ctx context.Context) (Role []domain.Role, err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "GetAllRoles",
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	return l.next.GetAllRoles(ctx)
}

func (l loggingService) DeleteRole(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "DeleteRole",
				"took", time.Since(begin),
				"err", err,
			)
		}
	}(time.Now())

	return l.next.DeleteRole(ctx, id)
}

func NewLoggingMiddleware(logger log.Logger, s service.RoleService) service.RoleService {
	return &loggingService{logger, s}
}
