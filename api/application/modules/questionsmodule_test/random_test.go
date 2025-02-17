package questionsmodule_test

import (
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/questionsmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
	"testing"

	"github.com/shelakel/go-ioc"
)

func TestGetRandomQuestions(t *testing.T) {
	c := ioc.NewContainer()
	c.MustRegister(func(f ioc.Factory) (interface{}, error) { return timemodule.NewClock(), nil }, (*timemodule.Clock)(nil), ioc.PerContainer)
	services := common.IocContainer(c)

	questionSet := questionsmodule.NewQuestionSetModel(
		modelmodule.NewModel(services),
		"",
		"",
		usersmodule.NewUser(
			modelmodule.NewModel(services),
			usersmodule.UserName("example"),
			"example.png",
			"",
			common.NewHasher(),
		),
	)

	for range make([]bool, 1000) {
		questionsmodule.NewQuestion(
			modelmodule.NewModel(services),
			&questionSet,
			"",
			"",
			nil,
		)
	}

	repeats := 5
	length := 100
	r1 := questionSet.GetRandomQuestions(services, length)
	idencital := true
	for idencital && repeats > 0 {
		repeats--
		r2 := questionSet.GetRandomQuestions(services, length)
		if len(r1) != length || len(r2) != length {
			t.Errorf("GetRandomQuestions returns invalid amount of questions")
			return
		}

		for i := range r1 {
			if r1[i].Id != r2[i].Id {
				idencital = false
				break
			}
		}
	}

	if idencital {
		t.Errorf("GetRandomQuestions returned identical sets")
		return
	}
}
