package middleware

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/rafaceo/go-test-auth/roles/domain"
	"github.com/rafaceo/go-test-auth/roles/service"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	requestError   metrics.Counter
	next           service.RoleService
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, counterE metrics.Counter, s service.RoleService) *instrumentingService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		requestError:   counterE,
		next:           s,
	}
}

func (s *instrumentingService) AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "AddRole"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AddRole(ctx, roleName, roleNameRu, notes, rights)
}

func (s *instrumentingService) EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditRole"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditRole(ctx, roleID, roleName, roleNameRu, notes, rights)
}

func (s *instrumentingService) GetRoles(ctx context.Context) (roles []domain.Role, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetRoles"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRoles(ctx)
}

func (s *instrumentingService) GetRoleRights(ctx context.Context, roleID int) (rights map[string][]string, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetRoleRights"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRoleRights(ctx, roleID)
}

func (s *instrumentingService) DeleteRole(ctx context.Context, roleID int) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "DeleteRole"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteRole(ctx, roleID)
}

func (s *instrumentingService) AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "AssignRoleToUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AssignRoleToUser(ctx, userID, roleID, merge)
}
