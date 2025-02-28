package middleware

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
	"github.com/rafaceo/go-test-auth/user_contexts/service"

	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	requestError   metrics.Counter
	next           service.UserContextService
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, counterE metrics.Counter, s service.UserContextService) service.UserContextService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		requestError:   counterE,
		next:           s,
	}
}

func (s *instrumentingService) AddUserContext(ctx context.Context, userID, merchantID uuid.UUID, global bool) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "AddUserContext"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AddUserContext(ctx, userID, merchantID, global)
}

func (s *instrumentingService) EditUserContext(ctx context.Context, userID uuid.UUID, global bool) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditUserContext"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditUserContext(ctx, userID, global)
}

func (s *instrumentingService) GetUserContexts(ctx context.Context, userID uuid.UUID) (contexts []domain.UserContext, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetUserContexts"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetUserContexts(ctx, userID)
}

func (s *instrumentingService) DeleteUserContext(ctx context.Context, userID, merchantID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "DeleteUserContext"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteUserContext(ctx, userID, merchantID)
}

func (s *instrumentingService) DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "DeleteAllUserContexts"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteAllUserContexts(ctx, userID)
}
