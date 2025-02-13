package matchmodule

import (
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/questionsmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
	"time"
)

// match dto
// match dto with relations
// match course dto
// match course dto with relations
// answered question dto
// player dto

// match dto

type MatchDto struct {
	modelmodule.ModelDto
	State           MatchState          `json:"state"`
	QuestionSetId   modelmodule.ModelId `json:"question_set_id"`
	QuestionsAmount int                 `json:"questions_amount"`
	HostUserId      modelmodule.ModelId `json:"host_user_id"`
}

func (model *MatchModel) Dto() MatchDto {
	var questionSetId modelmodule.ModelId
	if model.QuestionSetId != nil {
		questionSetId = *model.QuestionSetId
	}
	return MatchDto{
		ModelDto:        model.Model.Dto(),
		State:           model.State,
		QuestionSetId:   questionSetId,
		QuestionsAmount: model.QuestionsAmount,
		HostUserId:      model.HostUserId,
	}
}

// match dto with relations

type FullMatchDto struct {
	MatchDto
	// QuestionSet *questionsmodule.QuestionSetDto `json:"question_set"`
	Course  *FullMatchCourseDto `json:"course"`
	Players []PlayerDto         `json:"players"`
}

func (model *MatchModel) FullDto(c common.Ioc) FullMatchDto {
	// var questionSet *questionsmodule.QuestionSetDto
	// if model.QuestionSet != nil {
	// 	dto := model.QuestionSet.AnswerlessDto()
	// 	questionSet = &dto
	// }
	var course *FullMatchCourseDto
	if model.Course != nil {
		dto := model.Course.FullDto(c)
		course = &dto
	}
	var players []PlayerDto = []PlayerDto{}
	for _, player := range model.Players {
		players = append(players, player.Dto())
	}

	return FullMatchDto{
		MatchDto: model.Dto(),
		// QuestionSet: questionSet,
		Course:  course,
		Players: players,
	}
}

// match course dto

type MatchCourseDto struct {
	// modelmodule.ModelDto
	MatchId              modelmodule.ModelId          `json:"match_id"`
	CurrentQuestionIndex int                          `json:"current_question_index"`
	CurrentQuestion      *questionsmodule.QuestionDto `json:"current_question"`
	Step                 MatchCourseStep              `json:"step"`
	LastStep             string                       `json:"last_step"`
	NextStep             string                       `json:"next_step"`
}

func (model *MatchCourseModel) Dto(c common.Ioc) MatchCourseDto {
	var questionPtr *questionsmodule.QuestionDto
	if model.CurrentQuestion != -1 {
		question := model.questions()[model.CurrentQuestion].AnswerlessDto(c)
		questionPtr = &question
	}
	return MatchCourseDto{
		// ModelDto:             model.Model.Dto(),
		MatchId:              model.MatchId,
		CurrentQuestionIndex: model.CurrentQuestion,
		CurrentQuestion:      questionPtr,
		Step:                 model.Step,
		LastStep:             timemodule.FormatDate(time.Time(model.LastStep)),
		NextStep:             timemodule.FormatDate(time.Time(model.NextStep)),
	}
}

// match course dto with relations

type FullMatchCourseDto struct {
	MatchCourseDto
	AnsweredQuestions []AnsweredQuestionDto `json:"answered_questions"`
}

func (model *MatchCourseModel) FullDto(c common.Ioc) FullMatchCourseDto {
	var answeredQuestions []AnsweredQuestionDto = []AnsweredQuestionDto{}
	for _, answeredQuestion := range model.AnsweredQuestions {
		answeredQuestions = append(answeredQuestions, answeredQuestion.Dto())
	}
	return FullMatchCourseDto{
		MatchCourseDto:    model.Dto(c),
		AnsweredQuestions: answeredQuestions,
	}
}

// answered question dto

type AnsweredQuestionDto struct {
	modelmodule.ModelDto
	MatchCourseId     modelmodule.ModelId         `json:"match_course_id"`
	QuestionId        modelmodule.ModelId         `json:"question_id"`
	Question          questionsmodule.QuestionDto `json:"question"`
	AnsweredCorrectly bool                        `json:"answered_correctly"`
	AnswerTime        AnswerTime                  `json:"answer_time"`
	AnswredAt         AnsweredAt                  `json:"answered_at"`
	UserId            *modelmodule.ModelId        `json:"user_id"`
}

func (model *AnsweredQuestionModel) Dto() AnsweredQuestionDto {
	return AnsweredQuestionDto{
		ModelDto:          model.Model.Dto(),
		MatchCourseId:     model.MatchCourseId,
		QuestionId:        model.QuestionId,
		Question:          model.Question.Dto(),
		AnsweredCorrectly: model.AnsweredCorrectly,
		AnswerTime:        model.AnswerTime,
		AnswredAt:         model.AnswredAt,
		UserId:            model.UserId,
	}
}

// player dto

type PlayerDto struct {
	modelmodule.ModelDto
	MatchId modelmodule.ModelId `json:"match_id"`
	UserId  modelmodule.ModelId `json:"user_id"`
	User    usersmodule.UserDto `json:"user"`
	Online  bool                `json:"online"`
	Score   int                 `json:"score"`
}

func (model *PlayerModel) Dto() PlayerDto {
	return PlayerDto{
		ModelDto: model.Model.Dto(),
		MatchId:  model.MatchId,
		UserId:   model.UserId,
		User:     model.User.Dto(),
		Online:   model.Online,
		Score:    model.Score,
	}
}
