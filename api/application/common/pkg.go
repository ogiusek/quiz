package common

import (
	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

// we bind ourselfs to gorm, fasthttp and ioc.Container
// all of them can be changed but with each new package it is going to be harder
type Package interface {
	Db(db *gorm.DB)
	Services(c *ioc.Container)
	Variables(c Ioc) // mainly used by service groups
	Endpoints(r *router.Router, c IocScope)
}
