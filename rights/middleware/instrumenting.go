package middleware

import (
	"context"
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

func (s *instrumentingService) CreateUser(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "CreateUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateUser(ctx, name)
}

func (s *instrumentingService) EditUser(ctx context.Context, id uint, name, context string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditUser(ctx, id, name, context)
}

func (s *instrumentingService) GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (s *instrumentingService) EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (s *instrumentingService) RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) (err error) {
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

func (s *instrumentingService) GetUser(ctx context.Context, id uint) (name string, context string, rights map[string][]string, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetUser"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	name, context, rights, err = s.next.GetUser(ctx, id)
	return
}
