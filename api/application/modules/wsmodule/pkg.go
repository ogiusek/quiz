package wsmodule

import (
	"log"
	"quizapi/common"

	"github.com/fasthttp/router"
	"github.com/streadway/amqp"

	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

func DeclareExchanges(ch *amqp.Channel) {
	for exchange, queue := range binds {
		if err := ch.ExchangeDeclare(exchange, amqp.ExchangeDirect, false, false, false, false, nil); err != nil {
			log.Panic(err.Error())
		}

		if _, err := ch.QueueDeclare(queue, false, false, false, false, nil); err != nil {
			log.Panic(err.Error())
		}

		if err := ch.QueueBind(queue, "", exchange, false, nil); err != nil {
			log.Panic(err.Error())
		}
	}
}

type Package struct{}

func (Package) Db(db *gorm.DB) {
}

func AddSocketsStorageService(c *ioc.Container, storage socketInterface) {
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return storage, nil }, (*SocketStorage)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return storage, nil }, (*socketConnect)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewSocketsMessager(storage), nil }, (*SocketsMessager)(nil), ioc.PerContainer)
}

func (Package) Services(c *ioc.Container) {
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceGroup[Middleware](), nil }, (*common.ServiceGroup[Middleware])(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return common.NewServiceStorage[SocketId](), nil }, (*common.ServiceStorage[SocketId])(nil), ioc.PerScope)
}

func (Package) Variables(c common.Ioc) {
	var ch *amqp.Channel
	c.Inject(&ch)
	DeclareExchanges(ch)
}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
}

func HostApplication(c common.IocScope) {
	var sockets socketConnect
	c.Scope().Inject(&sockets)
	go sockets.Listen()
}
