package common

// agnostic, universal middleware
type Middleware interface {
	Handle(c Ioc, next func())
}

type middlewareImpl struct {
	handle func(c Ioc, next func())
}

func (m *middlewareImpl) Handle(c Ioc, next func()) {
	m.handle(c, next)
}

func NewMiddleware(handle func(c Ioc, next func())) Middleware {
	return &middlewareImpl{
		handle: handle,
	}
}

func RunMiddlewares(middlewares []Middleware, endpoint func(c Ioc), c Ioc) {
	var next func()
	next = func() {
		if len(middlewares) == 0 {
			endpoint(c)
			return
		}
		middleware := middlewares[0]
		middlewares = middlewares[1:]
		middleware.Handle(c, next)
	}

	next()
}
