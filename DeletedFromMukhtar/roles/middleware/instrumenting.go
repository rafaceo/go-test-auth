package middleware

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/service"
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

func (s *instrumentingService) AddRole(ctx context.Context, role domain.Role) (id uuid.UUID, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "AddRole"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AddRole(ctx, role)
}

func (s *instrumentingService) UpdateRole(ctx context.Context, role domain.Role) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditRole"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.UpdateRole(ctx, role)
}

func (s *instrumentingService) GetRoleByID(ctx context.Context, id uuid.UUID) (role *domain.Role, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetRoleById"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRoleByID(ctx, id)
}

func (s *instrumentingService) GetAllRoles(ctx context.Context) (right []domain.Role, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetAllRoles"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetAllRoles(ctx)
}

func (s *instrumentingService) DeleteRole(ctx context.Context, id uuid.UUID) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "DeleteRight"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteRole(ctx, id)
}
