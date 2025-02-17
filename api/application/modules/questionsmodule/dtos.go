package questionsmodule

import (
	"math/rand"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
)

// question set model dto
// question set model answerless dto
// question model dto
// question model answerless dto
// answer inputs answerless dto
// answer options answerless dto

// question set model dto

type QuestionSetDto struct {
	modelmodule.ModelDto
	Name        QuestionSetName        `json:"name"`
	Description QuestionSetDescription `json:"description"`
	Owner       usersmodule.UserDto    `json:"owner"`
	Questions   []QuestionDto          `json:"questions"`
}

func (model *QuestionSetModel) Dto() QuestionSetDto {
	var questions []QuestionDto
	for _, question := range model.Questions {
		questions = append(questions, question.Dto())
	}
	return QuestionSetDto{
		ModelDto:    model.Model.Dto(),
		Name:        model.Name,
		Description: model.Description,
		Owner:       model.Owner.Dto(),
		Questions:   questions,
	}
}

// question set model answerless dto

func (model *QuestionSetModel) AnswerlessDto(c common.Ioc) QuestionSetDto {
	var questions []QuestionDto
	for _, question := range model.Questions {
		questions = append(questions, question.AnswerlessDto(c))
	}
	return QuestionSetDto{
		ModelDto:    model.Model.Dto(),
		Name:        model.Name,
		Description: model.Description,
		Owner:       model.Owner.Dto(),
		Questions:   questions,
	}
}

// question model dto

type QuestionDto struct {
	modelmodule.ModelDto
	QuestionSetId modelmodule.ModelId `json:"question_set_id"`
	Question      Question            `json:"question"`
	AnswerType    AnswerType          `json:"answer_type"`
	Answer        any                 `json:"answer"`
}

func (model *QuestionModel) Dto() QuestionDto {
	return QuestionDto{
		ModelDto:      model.Model.Dto(),
		QuestionSetId: model.QuestionSetId,
		Question:      model.Question,
		AnswerType:    model.AnswerType,
		Answer:        model.GetAnswer(),
	}
}

// question model answerless dto

func (model *QuestionModel) AnswerlessDto(c common.Ioc) QuestionDto {
	return QuestionDto{
		ModelDto:      model.Model.Dto(),
		QuestionSetId: model.QuestionSetId,
		Question:      model.Question,
		AnswerType:    model.AnswerType,
		Answer:        model.GetAnswer().AnswerlessDto(c),
	}
}

// answer inputs answerless dto

func (answer *AnswerInputs) AnswerlessDto(c common.Ioc) any {
	return struct{}{}
}

// answer options answerless dto

func (answer *AnswerOptions) AnswerlessDto(c common.Ioc) any {
	var clock timemodule.Clock
	c.Inject(&clock)
	allAnswers := append(answer.Answers, answer.CorrectAnswer)
	answersLen := len(allAnswers)
	r := rand.New(rand.NewSource(clock.Now().UnixNano()))

	randomIndexes := r.Perm(answersLen)

	randomAnswers := make([]AnswerMessage, answersLen)
	for i, index := range randomIndexes {
		randomAnswers[i] = allAnswers[index]
	}
	return struct {
		Answers []AnswerMessage `json:"answers"`
	}{
		Answers: randomAnswers,
	}
}
