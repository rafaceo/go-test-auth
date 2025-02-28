package middleware

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/rafaceo/go-test-auth/rights/domain"
	"github.com/rafaceo/go-test-auth/rights/service"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	requestError   metrics.Counter
	next           service.RightsService
}

func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, counterE metrics.Counter, s service.RightsService) *instrumentingService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		requestError:   counterE,
		next:           s,
	}
}

// AddRight
func (s *instrumentingService) AddRights(ctx context.Context, module string, action []string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "AddRights"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AddRights(ctx, module, action)
}

// EditRight
func (s *instrumentingService) EditRight(ctx context.Context, id, module string, action []string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "EditRight"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.EditRight(ctx, id, module, action)
}

// GetAllRights
func (s *instrumentingService) GetAllRights(ctx context.Context) (rights []domain.Right, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetAllRights"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetAllRights(ctx)
}

// GetRightByName
func (s *instrumentingService) GetRightByName(ctx context.Context, module string) (right *domain.Right, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetRightByName"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRightByName(ctx, module)
}

// GetRightById
func (s *instrumentingService) GetRightById(ctx context.Context, id string) (right *domain.Right, err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "GetRightById"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRightById(ctx, id)
}

// DeleteRight
func (s *instrumentingService) DeleteRight(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		labels := []string{"method", "DeleteRight"}
		s.requestCount.With(labels...).Add(1)
		if err != nil {
			s.requestError.With(labels...).Add(1)
		}
		s.requestLatency.With(labels...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteRight(ctx, id)
}
