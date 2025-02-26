package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/rafaceo/go-test-auth/roles/domain"
	"github.com/rafaceo/go-test-auth/roles/service"
	"time"
)

type loggingService struct {
	logger log.Logger
	next   service.RoleService
}

func (l loggingService) AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "AddRole",
			"took", time.Since(begin),
			"roleName", roleName,
			"roleNameRu", roleNameRu,
			"notes", notes,
			"rights", rights,
			"err", err,
		)
	}(time.Now())

	return l.next.AddRole(ctx, roleName, roleNameRu, notes, rights)
}

func (l loggingService) EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "EditRole",
			"took", time.Since(begin),
			"roleID", roleID,
			"roleName", roleName,
			"roleNameRu", roleNameRu,
			"notes", notes,
			"rights", rights,
			"err", err,
		)
	}(time.Now())

	return l.next.EditRole(ctx, roleID, roleName, roleNameRu, notes, rights)
}

func (l loggingService) GetRoles(ctx context.Context) (roles []domain.Role, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetRoles",
			"took", time.Since(begin),
			"roles", roles,
			"err", err,
		)
	}(time.Now())

	return l.next.GetRoles(ctx)
}

func (l loggingService) GetRoleRights(ctx context.Context, roleID int) (rights map[string][]string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "GetRoleRights",
			"took", time.Since(begin),
			"roleID", roleID,
			"rights", rights,
			"err", err,
		)
	}(time.Now())

	return l.next.GetRoleRights(ctx, roleID)
}

func (l loggingService) DeleteRole(ctx context.Context, roleID int) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "DeleteRole",
			"took", time.Since(begin),
			"roleID", roleID,
			"err", err,
		)
	}(time.Now())

	return l.next.DeleteRole(ctx, roleID)
}

func (l loggingService) AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "AssignRoleToUser",
			"took", time.Since(begin),
			"userID", userID,
			"roleID", roleID,
			"merge", merge,
			"err", err,
		)
	}(time.Now())

	return l.next.AssignRoleToUser(ctx, userID, roleID, merge)
}

func NewLoggingMiddleware(logger log.Logger, s service.RoleService) service.RoleService {
	return &loggingService{logger, s}
}
