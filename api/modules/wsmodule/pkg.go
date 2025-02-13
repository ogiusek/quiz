package wsmodule

import (
	"quizapi/common"

	"github.com/fasthttp/router"

	"github.com/shelakel/go-ioc"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

type Package struct{}

func (Package) Db(db *gorm.DB) {
}

func (Package) Services(c *ioc.Container) {
	sockets := NewSockets()
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return sockets, nil }, (*SocketStorage)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return sockets, nil }, (*socketConnect)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewSocketsMessager(sockets), nil }, (*SocketsMessager)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceGroup[Middleware](), nil }, (*common.ServiceGroup[Middleware])(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[SocketId](), nil }, (*common.ServiceStorage[SocketId])(nil), ioc.PerScope)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[SocketConn](), nil }, (*common.ServiceStorage[SocketConn])(nil), ioc.PerScope)
}

func (Package) Variables(c common.Ioc) {}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
	var middlewaresGroup common.ServiceGroup[common.HttpMiddleware]
	scope := c.Scope()
	scope.Inject(&middlewaresGroup)
	middlewares := middlewaresGroup.GetAll()

	r.GET("/ws", func(ctx *fasthttp.RequestCtx) {
		common.RunHttpMiddlewares(middlewares, Connect, ctx, c.Scope())
	})

	var socketsMessager SocketsMessager
	scope.Inject(&socketsMessager)
}
