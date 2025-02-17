package main

import (
	"fmt"
	"log"
	"os"
	"quizapi/common"
	"quizapi/modules/usersmodule"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	c := ioc.NewContainer()

	// config

	config := GetValidConfig()
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return config, nil }, (*MainConfig)(nil), ioc.PerContainer)

	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewJwtConfig(config.JwtSecret), nil }, (*common.JwtConfig)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return config.UsersConfig, nil }, (*usersmodule.UserConfig)(nil), ioc.PerContainer)

	// log create logger

	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return logger, nil }, (*log.Logger)(nil), ioc.PerContainer)

	// config database with its logger

	var GormLogger gormLogger.Interface

	if config.Env == DevEnv {
		GormLogger = gormLogger.New(
			logger,
			gormLogger.Config{LogLevel: gormLogger.Info, Colorful: true},
		)
	}

	if config.Env == DeployEnv {
		GormLogger = gormLogger.Discard
	}

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{
		Logger: GormLogger,
	})
	if err != nil {
		log.Panic(err.Error())
	}

	// load packages

	for _, pkg := range packages {
		pkg.Db(db)
	}

	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return db, nil }, (*gorm.DB)(nil), ioc.PerContainer)
	for _, pkg := range packages {
		pkg.Services(c)
	}

	ioc := common.NewIocScope(func() common.Ioc { return common.IocContainer(c.Scope()) })
	for _, pkg := range packages {
		pkg.Variables(ioc.Scope())
	}

	r := router.New()
	for _, pkg := range packages {
		pkg.Endpoints(r, ioc)
	}

	log.Printf("starting server on :%d", config.Port)
	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", config.Port), func(ctx *fasthttp.RequestCtx) {
		// CORS
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type")

		// handle request
		r.Handler(ctx)
	}); err != nil {
		log.Panic(err.Error())
	}
}

// file storage				not written not injected
