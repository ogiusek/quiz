package questionsmodule

import (
	"errors"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
	"time"
)

// create question set
// change question set name
// change question set description
// search question sets
// get question set
// delete question set
// create question
// change question
// change question answer
// delete question

// create question set

type CreateQuestionSetArgs struct {
	Name        QuestionSetName        `json:"name"`
	Description QuestionSetDescription `json:"description"`
}

func (args *CreateQuestionSetArgs) Valid() []error {
	var errors []error
	for _, err := range args.Name.Valid() {
		errors = append(errors, common.ErrPath(err).Property("name"))
	}
	for _, err := range args.Description.Valid() {
		errors = append(errors, common.ErrPath(err).Property("description"))
	}
	return errors
}

func (args *CreateQuestionSetArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()

	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var userRepo usersmodule.UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return usersmodule.ErrUserNotFound
	}

	model := NewQuestionSetModel(
		modelmodule.NewModel(c),
		args.Name,
		args.Description,
		*user,
	)

	var questionSetRepo QuestionSetRepository
	c.Inject(&questionSetRepo)

	if err := questionSetRepo.Create(c, model); err != nil {
		return errQuestionSetConflict
	}

	var responseStorage common.ServiceStorage[common.Response]
	c.Inject(&responseStorage)
	responseStorage.Set(model.CreatedDto())

	return nil
}

// change question set name

type ChangeQuestionSetNameArgs struct {
	ModelId modelmodule.ModelId `json:"id"`
	NewName QuestionSetName     `json:"new_name"`
}

func (args *ChangeQuestionSetNameArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	for _, err := range args.NewName.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_name"))
	}
	return errors
}

func (args *ChangeQuestionSetNameArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()

	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionSetRepo QuestionSetRepository
	c.Inject(&questionSetRepo)

	set := questionSetRepo.GetById(c, args.ModelId)
	if set == nil {
		return ErrQuestionSetNotFound
	}

	set.ChangeName(args.NewName)

	if err := questionSetRepo.Update(c, *set); err == common.ErrRepositoryConflict {
		return errQuestionSetConflict
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// change question set description

type ChangeQuestionSetDescriptionArgs struct {
	ModelId        modelmodule.ModelId    `json:"id"`
	NewDescription QuestionSetDescription `json:"new_description"`
}

func (args *ChangeQuestionSetDescriptionArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	for _, err := range args.NewDescription.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_description"))
	}

	return errors
}

func (args *ChangeQuestionSetDescriptionArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionSetRepo QuestionSetRepository
	c.Inject(&questionSetRepo)

	set := questionSetRepo.GetById(c, args.ModelId)
	if set == nil {
		return ErrQuestionSetNotFound
	}

	set.ChangeDescription(args.NewDescription)

	if err := questionSetRepo.Update(c, *set); err == common.ErrRepositoryConflict {
		return errQuestionSetConflict
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// search question set

type SearchQuestionSetArgs SearchQuestionSets

type SearchQuestionSetRes struct {
	Found []QuestionSetDto `json:"found"`
	Time  string           `json:"when"`
}

// func (args *SearchQuestionSetArgs) Valid() []error { return args.Valid() }

func (args *SearchQuestionSetArgs) Handle(c common.Ioc) error {
	var repo QuestionSetRepository
	c.Inject(&repo)
	models := repo.Search(c, SearchQuestionSets(*args))
	var clock timemodule.Clock
	c.Inject(&clock)
	result := SearchQuestionSetRes{
		Found: make([]QuestionSetDto, 0),
		Time:  timemodule.FormatDate(clock.Now()),
	}

	var noDate time.Time
	if args.LastUpdate != noDate {
		result.Time = timemodule.FormatDate(args.LastUpdate)
	}

	for _, model := range models {
		result.Found = append(result.Found, model.Dto())
	}

	var responseStorage common.ServiceStorage[common.Response]
	c.Inject(&responseStorage)
	responseStorage.Set(result)

	return nil
}

// get question set

type GetQuestionSetArgs struct {
	ModelId modelmodule.ModelId `json:"id"`
}

type GetQuestionSetRes struct {
	Model QuestionSetDto `json:"model"`
}

func (args *GetQuestionSetArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	return errors
}

func (args *GetQuestionSetArgs) Handle(c common.Ioc) error {
	var repo QuestionSetRepository
	c.Inject(&repo)
	questionSet := repo.GetById(c, args.ModelId)
	if questionSet == nil {
		return ErrQuestionSetNotFound
	}
	var responseStorage common.ServiceStorage[common.Response]
	c.Inject(&responseStorage)
	responseStorage.Set(GetQuestionSetRes{
		Model: questionSet.Dto(),
	})
	return nil
}

// delete question set

type DeleteQuestionSetArgs struct {
	ModelId modelmodule.ModelId `json:"id"`
}

func (args *DeleteQuestionSetArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	return errors
}

func (args *DeleteQuestionSetArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionSetRepo QuestionSetRepository
	c.Inject(&questionSetRepo)
	questionSet := questionSetRepo.GetById(c, args.ModelId)
	if questionSet == nil {
		return ErrQuestionSetNotFound
	}
	if questionSet.OwnerId != session.UserId {
		return errQuestionSetForbidden
	}
	if err := questionSetRepo.Delete(c, questionSet.Id); err == common.ErrRepositoryNotFound {
		return ErrQuestionSetNotFound
	}
	return nil
}

// create question

type CreateQuestionArgs struct {
	QuestionSetId modelmodule.ModelId `json:"question_set_id"`
	Question      Question            `json:"question"`
	AnswerType    AnswerType          `json:"answer_type"`
	Answer        any                 `json:"answer"`
}

func (args *CreateQuestionArgs) Valid() []error {
	var errors []error
	for _, err := range args.QuestionSetId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("question_set_id"))
	}
	for _, err := range args.Question.Valid() {
		errors = append(errors, common.ErrPath(err).Property("question"))
	}
	for _, err := range args.AnswerType.Valid() {
		errors = append(errors, common.ErrPath(err).Property("answer_type"))
	}
	if answer := args.AnswerType.RawAnswer(args.Answer); answer != nil {
		for _, err := range answer.Valid() {
			errors = append(errors, common.ErrPath(err).Property("answer"))
		}
	}
	return errors
}

func (args *CreateQuestionArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionSetRepo QuestionSetRepository
	c.Inject(&questionSetRepo)
	quesionSet := questionSetRepo.GetById(c, args.QuestionSetId)
	if quesionSet == nil {
		return ErrQuestionSetNotFound
	}
	if quesionSet.OwnerId != session.UserId {
		return errQuestionSetForbidden
	}

	var questionRepo QuestionRepository
	c.Inject(&questionRepo)
	question := NewQuestion(
		modelmodule.NewModel(c),
		quesionSet,
		args.Question,
		args.AnswerType,
		args.AnswerType.RawAnswer(args.Answer),
	)
	if err := questionRepo.Create(c, question); err == common.ErrRepositoryConflict {
		return errQuestionConflict
	} else if err == common.ErrHttpParallelModification {
		return common.ErrHttpParallelModification
	}

	var responseStorage common.ServiceStorage[common.Response]
	c.Inject(&responseStorage)
	responseStorage.Set(question.CreatedDto())

	return nil
}

// change question

type ChangeQuestionArgs struct {
	ModelId     modelmodule.ModelId `json:"id"`
	NewQuestion Question            `json:"new_question"`
}

func (args *ChangeQuestionArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	for _, err := range args.NewQuestion.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_question"))
	}
	return errors
}

func (args *ChangeQuestionArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionRepo QuestionRepository
	c.Inject(&questionRepo)
	question := questionRepo.GetById(c, args.ModelId)
	if question == nil {
		return errQuestionNotFound
	}
	if question.QuestionSet.Owner.Id != session.UserId {
		return errQuestionSetForbidden
	}
	question.ChangeQuestion(args.NewQuestion)
	if err := questionRepo.Update(c, *question); err == common.ErrRepositoryConflict {
		return errQuestionConflict
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// change question answer

type ChangeQuestionAnswerArgs struct {
	ModelId       modelmodule.ModelId `json:"id"`
	NewAnswerType AnswerType          `json:"new_answer_type"`
	NewAnswer     any                 `json:"new_answer"`
}

var (
	errMissingField error = errors.New("missing field value")
)

func (args *ChangeQuestionAnswerArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	for _, err := range args.NewAnswerType.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_answer_type"))
	}
	if answer := args.NewAnswerType.RawAnswer(args.NewAnswer); answer != nil {
		for _, err := range answer.Valid() {
			errors = append(errors, common.ErrPath(err).Property("new_answer"))
		}
	} else {
		errors = append(errors, common.ErrPath(errMissingField).Property("new_answer"))
	}
	return errors
}

func (args *ChangeQuestionAnswerArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionRepo QuestionRepository
	c.Inject(&questionRepo)
	question := questionRepo.GetById(c, args.ModelId)
	if question == nil {
		return errQuestionNotFound
	}
	if question.QuestionSet.Owner.Id != session.UserId {
		return errQuestionSetForbidden
	}
	question.ChangeAnswer(args.NewAnswerType, args.NewAnswerType.RawAnswer(args.NewAnswer))
	if err := questionRepo.Update(c, *question); err == common.ErrRepositoryConflict {
		return errQuestionConflict
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// delete question

type DeleteQuestionArgs struct {
	ModelId modelmodule.ModelId `json:"id"`
}

func (args *DeleteQuestionArgs) Valid() []error {
	var errors []error
	for _, err := range args.ModelId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	return errors
}

func (args *DeleteQuestionArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var questionRepo QuestionRepository
	c.Inject(&questionRepo)
	question := questionRepo.GetById(c, args.ModelId)
	if question == nil {
		return errQuestionNotFound
	}
	if question.QuestionSet.OwnerId != session.UserId {
		return errQuestionSetForbidden
	}
	if err := questionRepo.Delete(c, question.Id); err == common.ErrRepositoryNotFound {
		return errQuestionNotFound
	}

	return nil
}
