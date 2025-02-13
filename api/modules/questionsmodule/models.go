package questionsmodule

import (
	"encoding/json"
	"math/rand"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
)

// question model
// question set model

// question model

type QuestionModel struct {
	modelmodule.Model
	QuestionSet   *QuestionSetModel   `gorm:"foreignKey:QuestionSetId"`
	QuestionSetId modelmodule.ModelId `gorm:"column:question_set_id;not null"`
	Question      Question            `gorm:"column:question;not null"`
	AnswerType    AnswerType          `gorm:"column:answer_type;not null"`
	Answer        []byte              `gorm:"column:answer;type:jsonb;not null"`
}

func (QuestionModel) TableName() string { return "questions" }

func NewQuestion(model modelmodule.Model, questionSet *QuestionSetModel, question Question, answerType AnswerType, answer Answer) QuestionModel {
	encodedAnswer, _ := json.Marshal(answer)
	questionModel := QuestionModel{
		Model:         model,
		QuestionSet:   questionSet,
		QuestionSetId: questionSet.Id,
		Question:      question,
		AnswerType:    answerType,
		Answer:        encodedAnswer,
	}
	questionSet.Questions = append(questionSet.Questions, questionModel)
	return questionModel
}

func (model *QuestionModel) GetAnswer() Answer {
	return model.AnswerType.Answer(model.Answer)
}

func (model *QuestionModel) ChangeQuestion(question Question) {
	model.Question = question
}

func (model *QuestionModel) ChangeAnswer(answerType AnswerType, answer Answer) {
	model.AnswerType = answerType
	model.Answer, _ = json.Marshal(answer)
}

// question set model

type QuestionSetModel struct {
	modelmodule.Model
	Name        QuestionSetName        `gorm:"column:name;not null"`
	Description QuestionSetDescription `gorm:"column:description;not null"`
	OwnerId     modelmodule.ModelId    `gorm:"column:owner_id;constraint:OnDelete:CASCADE;not null"`
	Owner       *usersmodule.UserModel `gorm:"foreignKey:OwnerId"`
	Questions   []QuestionModel        `gorm:"foreignKey:QuestionSetId;constraint:OnDelete:CASCADE"`
}

func (QuestionSetModel) TableName() string { return "question_sets" }

func NewQuestionSetModel(model modelmodule.Model, name QuestionSetName, description QuestionSetDescription, owner usersmodule.UserModel) QuestionSetModel {
	return QuestionSetModel{
		Model:       model,
		Name:        name,
		Description: description,
		OwnerId:     owner.Id,
		Owner:       &owner,
		Questions:   make([]QuestionModel, 0),
	}
}

func (model *QuestionSetModel) ChangeName(name QuestionSetName) {
	model.Name = name
}

func (model *QuestionSetModel) ChangeDescription(description QuestionSetDescription) {
	model.Description = description
}

func (model *QuestionSetModel) ChangeOwner(user usersmodule.UserModel) {
	model.Owner = &user
	model.OwnerId = user.Id
}

func (model *QuestionSetModel) GetRandomQuestions(c common.Ioc, questionsCount int) []*QuestionModel {
	if len(model.Questions) < questionsCount {
		return make([]*QuestionModel, 0)
	}
	var clock timemodule.Clock
	c.Inject(&clock)

	// New recommended way to generate random numbers
	r := rand.New(rand.NewSource(clock.Now().UnixNano()))

	// Get 10 unique random indices
	randomIndexes := r.Perm(len(model.Questions))[:questionsCount]

	questions := make([]*QuestionModel, questionsCount)
	for i, index := range randomIndexes {
		questions[i] = &model.Questions[index]
	}

	return questions
}
