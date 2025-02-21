package middleware

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/rafaceo/go-test-auth/rights/service"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	requestError   metrics.Counter
	next           service.UserService
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, counterE metrics.Counter, s service.UserService) service.UserService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		requestError:   counterE,
		next:           s,
	}
}

func (s *instrumentingService) CreateUser(ctx context.Context, phone string, passwordHash string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "CreateUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateUser(ctx, phone, passwordHash)
}

func (s *instrumentingService) EditUser(ctx context.Context, id uuid.UUID, phone, password string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditUser(ctx, id, phone, password)
}

func (s *instrumentingService) GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GrantRightsToUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GrantRightsToUser(ctx, id, rights)
}

func (s *instrumentingService) EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditRightsToUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditRightsToUser(ctx, id, rights)
}

func (s *instrumentingService) RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "RevokeRightsFromUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.RevokeRightsFromUser(ctx, id, rights)
}

func (s *instrumentingService) GetUser(ctx context.Context, id uuid.UUID) (phone string, password string, created_at string, updated_at string, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	phone, password, created_at, updated_at, err = s.next.GetUser(ctx, id)
	return
}

func (s *instrumentingService) GetUserRights(ctx context.Context, id uuid.UUID) (rights map[string][]string, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetUserRights"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rights, err = s.next.GetUserRights(ctx, id)
	return
}
