package main

import (
	"log"
	"quizapi/common"
	"quizapi/modules/wsmodule"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

type MainPackage struct{}

func (MainPackage) Db(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		log.Panic(err.Error())
	}

	sqlDb.SetMaxIdleConns(8)
	sqlDb.SetMaxOpenConns(16)
}

func (MainPackage) Services(c *ioc.Container) {
	var config MainConfig
	c.MustResolve(&config)

	conn, err := amqp.Dial(config.RabbitMqUrl)
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err.Error())
	}
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return ch, nil }, (**amqp.Channel)(nil), ioc.PerContainer)

	// handling requests services
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceGroup[common.Middleware](), nil }, (*common.ServiceGroup[common.Middleware])(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceGroup[common.HttpMiddleware](), nil }, (*common.ServiceGroup[common.HttpMiddleware])(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[common.Response](), nil }, (*common.ServiceStorage[common.Response])(nil), ioc.PerScope)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[error](), nil }, (*common.ServiceStorage[error])(nil), ioc.PerScope)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[*gorm.DB](), nil }, (*common.ServiceStorage[*gorm.DB])(nil), ioc.PerScope)

	// application services
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewHasher(), nil }, (*common.Hasher)(nil), ioc.PerContainer)
}

func (MainPackage) Variables(c common.Ioc) {
	var generalMiddlewareGroup common.ServiceGroup[common.Middleware]
	c.Inject(&generalMiddlewareGroup)

	var httpMiddlewaresGroup common.ServiceGroup[common.HttpMiddleware]
	c.Inject(&httpMiddlewaresGroup)

	var wsMiddlewaresGroup common.ServiceGroup[wsmodule.Middleware]
	c.Inject(&wsMiddlewaresGroup)

	var db *gorm.DB
	c.Inject(&db)

	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	dbStorage.Set(db)

	generalMiddleware := func(c common.Ioc, next func()) {
		db.Transaction(func(tx *gorm.DB) error {
			var dbStorage common.ServiceStorage[*gorm.DB]
			c.Inject(&dbStorage)
			dbStorage.Set(tx)
			next()
			var errStorage common.ServiceStorage[error]
			c.Inject(&errStorage)
			errPtr := errStorage.Get()
			if errPtr == nil {
				return nil
			}

			return *errPtr
		})
	}

	generalMiddlewareGroup.Add(common.NewMiddleware(generalMiddleware))
	httpMiddlewaresGroup.Add(common.NewHttpMiddleware(func(ctx *fasthttp.RequestCtx, c common.Ioc, next func()) { generalMiddleware(c, next) }))
	wsMiddlewaresGroup.Add(wsmodule.NewMiddleware(func(c common.Ioc, message []byte, next func()) { generalMiddleware(c, next) }))
}

func (MainPackage) Endpoints(r *router.Router, c common.IocScope) {}
