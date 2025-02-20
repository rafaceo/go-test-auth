package instrumenting

type SimpleFactory struct{}

func (s *SimpleFactory) CreateInstrumentingService(subsystem string) *Service {
	counter, duration, counterError := GetMetricsBySubsystem(subsystem)

	return &Service{
		requestCount:   counter,
		requestLatency: duration,
		requestError:   counterError,
	}
}
