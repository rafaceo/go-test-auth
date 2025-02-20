package instrumenting

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type Service struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	requestError   metrics.Counter
}

func (s *Service) SaveMetricsByMethodName(methodName string, err error) {
	begin := time.Now()
	labelValues := []string{"method", methodName}
	s.requestCount.With(labelValues...).Add(1)
	if err != nil {
		s.requestError.With(labelValues...).Add(1)
	}
	s.requestLatency.With(labelValues...).Observe(time.Since(begin).Seconds())
}
