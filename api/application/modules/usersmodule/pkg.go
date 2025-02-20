package usersmodule

import (
	"quizapi/common"
	"quizapi/modules/wsmodule"
	"time"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

type Package struct{}

func (Package) Db(db *gorm.DB) {
	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&UserSocket{})
}

func (Package) Services(c *ioc.Container) {
	var db *gorm.DB
	c.MustResolve(&db)
	// repos
	userSocketRepo := NewUserSocketRepository()
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewUserRepository(), nil }, (*UserRepository)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return userSocketRepo, nil }, (*UserSocketRepo)(nil), ioc.PerScope)

	// session
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[SessionDto](), nil }, (*common.ServiceStorage[SessionDto])(nil), ioc.PerScope)

	wsmodule.AddSocketsStorageService(c, wsmodule.NewSocketStorage(common.IocContainer(c), func(socketId wsmodule.SocketId, meta []byte) bool {
		services := common.IocContainer(c)

		var generalMiddlewareGroup common.ServiceGroup[common.Middleware]
		services.Inject(&generalMiddlewareGroup)
		middlewares := generalMiddlewareGroup.GetAll()

		var canConnect bool
		common.RunMiddlewares(middlewares, func(c common.Ioc) {
			session := SessionDto{}
			if err := session.Decode(services, string(meta)); err != nil {
				canConnect = false
				return
			}
			userSocketRepo.CreateSocket(services, NewUserSocket(socketId, session.UserId))
			canConnect = true
		}, services)

		return canConnect
	}))
}

func (Package) Variables(c common.Ioc) {
	var httpMiddlewares common.ServiceGroup[common.HttpMiddleware]
	c.Inject(&httpMiddlewares)
	httpMiddlewares.Add(common.NewHttpMiddleware(authorizeHttpMiddleware))

	var wsMiddlewares common.ServiceGroup[wsmodule.Middleware]
	c.Inject(&wsMiddlewares)
	wsMiddlewares.Add(wsmodule.NewMiddleware(authorizeWsMiddleware))

	var storage wsmodule.SocketStorage
	c.Inject(&storage)
	var middlewareGroup common.ServiceGroup[common.Middleware]
	middlewares := middlewareGroup.GetAll()
	c.Inject(&middlewareGroup)
	storage.OnStart(func() {
		common.RunMiddlewares(middlewares, func(c common.Ioc) {
			var db *gorm.DB
			c.Inject(&db)
			db.Where("1 = 1").Delete(&UserSocket{})
		}, c)
	})
	storage.OnClose(func(id wsmodule.SocketId) {
		common.RunMiddlewares(middlewares, func(c common.Ioc) {
			go func() {
				// TODO find way to order listeners
				// this is a temporary solution highly unstable
				// in match module there is on close listener which should be executed first
				time.Sleep(time.Millisecond)
				var db *gorm.DB
				c.Inject(&db)
				db.Where("socket_id = ?", id).Delete(&UserSocket{})
			}()
		}, c)
	})
}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
	var socketStorage wsmodule.SocketStorage
	var middlewareGroup common.ServiceGroup[common.Middleware]
	c.Scope().Inject(&socketStorage)
	c.Scope().Inject(&middlewareGroup)
	middlewares := middlewareGroup.GetAll()

	socketStorage.OnStart(func() {
		common.RunMiddlewares(middlewares, func(c common.Ioc) {
			var dbStorage common.ServiceStorage[*gorm.DB]
			c.Inject(&dbStorage)
			dbStorage.MustGet().Where("1 = 1").Delete(&UserSocket{})
		}, c.Scope())
	})

	r.POST("/api/user/register", common.HttpEndpoint(c, common.FromBody, (*RegisterArgs)(nil)))
	r.POST("/api/user/log-in", common.HttpEndpoint(c, common.FromBody, (*LogInArgs)(nil)))
	r.POST("/api/user/refresh", common.HttpEndpoint(c, common.FromBody, (*RefreshArgs)(nil)))
	r.GET("/api/user/profile", common.HttpEndpoint(c, common.FromQuery, (*ProfileArgs)(nil)))
	r.POST("/api/user/change-name", common.HttpEndpoint(c, common.FromBody, (*ChangeNameArgs)(nil)))
	// r.POST("/api/user/change-profile-picture", common.HttpEndpoint(c, common.FromBody, (*ChangeProfilePictureArgs)(nil)))
	r.POST("/api/user/change-password", common.HttpEndpoint(c, common.FromBody, (*ChangePasswordArgs)(nil)))
}
