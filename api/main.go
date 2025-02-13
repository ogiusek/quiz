package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"quizapi/common"
	"strings"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetValidConfig() MainConfig {
	configPath := flag.String("config", "./config.json", "Path to the configuration file")
	flag.Parse()
	configData := *configPath

	file, err := os.Open(configData)
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()

	var config MainConfig
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&config); err != nil {
		log.Panic(err.Error())
	}

	if errors := config.Valid(); len(errors) != 0 {
		var messages []string
		for _, err := range errors {
			messages = append(messages, err.Error())
		}
		log.Panicf("errors:\n```\n%s\n```", strings.Join(messages, "\n"))
	}

	return config
}

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
	{
		config := GetValidConfig()
		c.MustRegister(func(f ioc.Factory) (interface{}, error) { return config, nil }, (*MainConfig)(nil), ioc.PerContainer)
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
