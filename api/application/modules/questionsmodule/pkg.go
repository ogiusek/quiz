package questionsmodule

import (
	"quizapi/common"

	"github.com/fasthttp/router"
	"github.com/shelakel/go-ioc"
	"gorm.io/gorm"
)

type Package struct{}

func (Package) Db(db *gorm.DB) {
	db.AutoMigrate(&QuestionModel{})
	db.AutoMigrate(&QuestionSetModel{})
}

func (Package) Services(c *ioc.Container) {
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewQuestionSetRepository(), nil }, (*QuestionSetRepository)(nil), ioc.PerContainer)
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return NewQuestionRepository(), nil }, (*QuestionRepository)(nil), ioc.PerContainer)
}

func (Package) Variables(c common.Ioc) {
}

func (Package) Endpoints(r *router.Router, c common.IocScope) {
	r.POST("/api/question-set/create", common.HttpEndpoint(c, common.FromBody, (*CreateQuestionSetArgs)(nil)))
	r.POST("/api/question-set/change-name", common.HttpEndpoint(c, common.FromBody, (*ChangeQuestionSetNameArgs)(nil)))
	r.POST("/api/question-set/change-description", common.HttpEndpoint(c, common.FromBody, (*ChangeQuestionSetDescriptionArgs)(nil)))
	r.GET("/api/question-set/search", common.HttpEndpoint(c, common.FromQuery, (*SearchQuestionSetArgs)(nil)))
	r.GET("/api/question-set/get", common.HttpEndpoint(c, common.FromQuery, (*GetQuestionSetArgs)(nil)))
	r.DELETE("/api/question-set/delete", common.HttpEndpoint(c, common.FromBody, (*DeleteQuestionSetArgs)(nil)))
	r.POST("/api/question/create", common.HttpEndpoint(c, common.FromBody, (*CreateQuestionArgs)(nil)))
	r.POST("/api/question/change-question", common.HttpEndpoint(c, common.FromBody, (*ChangeQuestionArgs)(nil)))
	r.POST("/api/question/change-answer", common.HttpEndpoint(c, common.FromBody, (*ChangeQuestionAnswerArgs)(nil)))
	r.DELETE("/api/question/delete", common.HttpEndpoint(c, common.FromBody, (*DeleteQuestionArgs)(nil)))
}
