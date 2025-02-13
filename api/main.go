package main

import (
	"log"
	"quizapi/common"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "host=localhost port=5432 user=username password=password dbname=database"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		log.Panic(err.Error())
	}
	for _, pkg := range packages {
		pkg.Db(db)
	}

	c := ioc.NewContainer()
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

	log.Print("started server")
	if err := fasthttp.ListenAndServe(":5050", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type")
		r.Handler(ctx)
	}); err != nil {
		log.Panic(err.Error())
	}
}

// file storage				not written not injected
