package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	service "github.com/rafaceo/go-test-auth/rights/service"
)

type loggingService struct {
	next   service.UserService
	logger log.Logger
}

func (l *loggingService) CreateUser(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			_ = l.logger.Log(
				"method", "CreateUser",
				"took", time.Since(begin),
				"name", name,
				"err", err,
			)
		}
	}(time.Now())

	return l.next.CreateUser(ctx, name)
}
