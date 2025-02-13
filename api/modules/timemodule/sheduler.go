package timemodule

import (
	"quizapi/common"
	"time"
)

type Scheduler interface {
	Schedule(c common.Ioc, t time.Time, method func())
}

type schedulerImpl struct {
	injected    bool
	middlewares []common.Middleware
}

func (s *schedulerImpl) RunEnpoint(c common.Ioc, endpoint func(c common.Ioc)) {
	if !s.injected {
		var middlewareGroup common.ServiceGroup[common.Middleware]
		c.Inject(&middlewareGroup)
		s.middlewares = middlewareGroup.GetAll()
		s.injected = true
	}

	common.RunMiddlewares(s.middlewares, endpoint, c)
}

func (s *schedulerImpl) Schedule(c common.Ioc, t time.Time, method func()) {
	time.AfterFunc(time.Until(t), func() {
		s.RunEnpoint(c, func(c common.Ioc) { method() })
	})
}

func NewScheduler() Scheduler {
	return &schedulerImpl{
		injected: false,
	}
}
