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
}

func (Package) Services(c *ioc.Container) {
	var db *gorm.DB
	c.MustResolve(&db)
	// repos
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewUserRepository(), nil }, (*UserRepository)(nil), ioc.PerContainer)

	// session
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[SessionDto](), nil }, (*common.ServiceStorage[SessionDto])(nil), ioc.PerScope)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) {
		return UserConfig{
			SessionTokenExpirationTime: time.Minute * 5, // * 1024 * 1024, // (TODO DEVELOPMENT) remove comment when developing api and add it when finished
			RefreshTokenExpirationTime: time.Hour * 24 * 365,
		}, nil
	}, (*UserConfig)(nil), ioc.PerContainer)
}

func (Package) Variables(c common.Ioc) {
	var httpMiddlewares common.ServiceGroup[common.HttpMiddleware]
	c.Inject(&httpMiddlewares)
	httpMiddlewares.Add(common.NewHttpMiddleware(authorizeHttpMiddleware))

	var wsMiddlewares common.ServiceGroup[wsmodule.Middleware]
	c.Inject(&wsMiddlewares)
	wsMiddlewares.Add(wsmodule.NewMiddleware(authorizeWsMiddleware))
}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
	r.POST("/api/user/register", common.HttpEndpoint(c, common.FromBody, (*RegisterArgs)(nil)))
	r.POST("/api/user/log-in", common.HttpEndpoint(c, common.FromBody, (*LogInArgs)(nil)))
	r.POST("/api/user/refresh", common.HttpEndpoint(c, common.FromBody, (*RefreshArgs)(nil)))
	r.GET("/api/user/profile", common.HttpEndpoint(c, common.FromQuery, (*ProfileArgs)(nil)))
	r.POST("/api/user/change-name", common.HttpEndpoint(c, common.FromBody, (*ChangeNameArgs)(nil)))
	// r.POST("/api/user/change-profile-picture", common.HttpEndpoint(c, common.FromBody, (*ChangeProfilePictureArgs)(nil)))
	r.POST("/api/user/change-password", common.HttpEndpoint(c, common.FromBody, (*ChangePasswordArgs)(nil)))
}
