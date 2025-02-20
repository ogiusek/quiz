package timemodule

import (
	"quizapi/common"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

type Package struct {
}

func (Package) Db(db *gorm.DB) {}

func (Package) Services(c *ioc.Container) {
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewClock(), nil }, (*Clock)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewScheduler(), nil }, (*Scheduler)(nil), ioc.PerContainer)
}

func (Package) Variables(c common.Ioc)                        {}
func (Package) Endpoints(r *router.Router, c common.IocScope) {}
