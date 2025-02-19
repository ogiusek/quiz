package matchmodule

import (
	"quizapi/common"
	"quizapi/modules/eventsmodule"
	"quizapi/modules/wsmodule"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

type Package struct{}

func (Package) Db(db *gorm.DB) {
	db.AutoMigrate(
		&MatchModel{},
		&MatchCourseModel{},
		&AnsweredQuestionModel{},
		&PlayerModel{},
		&MatchCourseQuestionModel{},
	)
	db.SetupJoinTable(&MatchCourseModel{}, "Questions", &MatchCourseQuestionModel{})
}

func (Package) Services(c *ioc.Container) {
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewMatchRepository(), nil }, (*MatchRepository)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewMatchCourseRepository(), nil }, (*MatchCourseRepository)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewAnsweredQuestionsRepository(), nil }, (*AnsweredQuestionsRepository)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewPlayerRepository(), nil }, (*PlayerRepository)(nil), ioc.PerContainer)
}

func (Package) Variables(c common.Ioc) {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	for _, event := range events {
		eventManager.Reserve(event.Topic)
	}

	for topic, handler := range eventHandlers {
		eventManager.Listen(topic, handler)
	}
}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
	var socketsMessager wsmodule.SocketsMessager
	c.Scope().Inject(&socketsMessager)

	socketsMessager.Listen("match/host", wsmodule.WsEndpoint(c, (*HostArgs)(nil)))
	socketsMessager.Listen("match/active", wsmodule.WsEndpoint(c, (*ActiveGameArgs)(nil)))
	socketsMessager.Listen("match/join", wsmodule.WsEndpoint(c, (*JoinArgs)(nil)))
	socketsMessager.Listen("match/quit", wsmodule.WsEndpoint(c, (*QuitArgs)(nil)))
	socketsMessager.Listen("match/change-question-set", wsmodule.WsEndpoint(c, (*ChangeQuestionSetArgs)(nil)))
	socketsMessager.Listen("match/change-questions-amount", wsmodule.WsEndpoint(c, (*ChangeQuestionsAmountArgs)(nil)))
	socketsMessager.Listen("match/start", wsmodule.WsEndpoint(c, (*StartArgs)(nil)))
	socketsMessager.Listen("match/reset", wsmodule.WsEndpoint(c, (*ResetArgs)(nil)))
	socketsMessager.Listen("match/answer", wsmodule.WsEndpoint(c, (*AnswerArgs)(nil)))

	var sockets wsmodule.SocketStorage
	c.Scope().Inject(&sockets)
	scope := c.Scope()

	var middlewares common.ServiceGroup[common.Middleware]
	scope.Inject(&middlewares)
	sockets.OnStart(func() {
		common.RunMiddlewares(middlewares.GetAll(), func(c common.Ioc) {
			var dbStorage common.ServiceStorage[*gorm.DB]
			c.Inject(&dbStorage)
			dbStorage.MustGet().Where("1 = 1").Delete(&MatchModel{})
		}, c.Scope())
	})

	activeGame := wsmodule.WsEndpoint(c, (*ActiveGameArgs)(nil))
	sockets.OnConnect(func(id wsmodule.SocketId) { activeGame(id, wsmodule.EmptyPayload) })

	quit := wsmodule.WsEndpoint(c, (*QuitArgs)(nil))
	sockets.OnClose(func(id wsmodule.SocketId) { quit(id, wsmodule.EmptyPayload) })
}
