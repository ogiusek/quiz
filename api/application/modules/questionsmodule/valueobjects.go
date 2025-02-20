package questionsmodule

import (
	"encoding/json"
	"errors"
	"fmt"
	"quizapi/common"
	"strings"
)

// answer message
// answer
// answer options
// answer input
// answer inputs
// answer type
// question content
// question set name
// question set description

// answer message

type AnswerMessage string

func (AnswerMessage) GormDataType() string { return "varchar(64)" }

func (answer *AnswerMessage) Valid() []error {
	var errs []error
	if *answer == "" {
		errs = append(errs, errors.New("answer cannot be empty"))
	}
	if len(*answer) > 64 {
		errs = append(errs, errors.New("name cannot have more than 64 characters"))
	}
	return errs
}

// answer

type Answer interface {
	common.Validable
	IsCorrectAnswer(AnswerMessage) bool
	AnswerlessDto(common.Ioc) any
}

// answer options

var AnswerOptionsAnswerType AnswerType = "o"

type AnswerOptions struct {
	Answers       []AnswerMessage `json:"answers"`
	CorrectAnswer AnswerMessage   `json:"correct_answer"`
}

var (
	errToFewOptionsAnswers error = errors.New("there has to be 1 or 3 other answers")
)

func (answer *AnswerOptions) Valid() []error {
	var errors []error
	if len(answer.Answers) != 1 && len(answer.Answers) != 3 {
		errors = append(errors, common.ErrPath(errToFewOptionsAnswers).Property("answers[]"))
	}
	for i, answer := range answer.Answers {
		for _, err := range answer.Valid() {
			errors = append(errors, common.ErrPath(err).Property(fmt.Sprintf("answers[%d]", i)))
		}
	}
	for _, err := range answer.CorrectAnswer.Valid() {
		errors = append(errors, common.ErrPath(err).Property("correct_answer"))
	}
	return errors
}

func (answer *AnswerOptions) IsCorrectAnswer(message AnswerMessage) bool {
	return answer.CorrectAnswer == message
}

func NewAnswerOptions(correctAnswer AnswerMessage, otherAnswers []AnswerMessage) Answer {
	return &AnswerOptions{
		Answers:       otherAnswers,
		CorrectAnswer: correctAnswer,
	}
}

// answer input

// var answerInputTrim string = " \n"

type AnswerInput struct {
	Answer        AnswerMessage `json:"answer"`
	CaseSensitive bool          `json:"case_sensitive"`
}

func (answer *AnswerInput) Valid() []error {
	var errors []error
	for _, err := range answer.Answer.Valid() {
		errors = append(errors, common.ErrPath(err).Property("answer"))
	}
	return errors
}

func (input *AnswerInput) Matches(message AnswerMessage) bool {
	// message = AnswerMessage(strings.Trim(string(message), answerInputTrim))
	if input.CaseSensitive {
		return message == input.Answer
	}
	return strings.EqualFold(string(message), string(input.Answer))
}

// answer inputs

var AnswerInputsAnswerType AnswerType = "i"

type AnswerInputs struct {
	CorrectAnswers []AnswerInput `json:"correct_answers"`
}

var (
	errToFewInputAnswers error = errors.New("there have to be at least 1 correct input answer")
)

func (answer *AnswerInputs) Valid() []error {
	var errors []error
	if len(answer.CorrectAnswers) == 0 {
		errors = append(errors, common.ErrPath(errToFewInputAnswers).Property("correct_answers[]"))
	}
	for i, input := range answer.CorrectAnswers {
		for _, err := range input.Valid() {
			errors = append(errors, common.ErrPath(err).Property(fmt.Sprintf("correct_answers[%d]", i)))
		}
	}
	return errors
}

func (answer *AnswerInputs) IsCorrectAnswer(message AnswerMessage) bool {
	for _, answer := range answer.CorrectAnswers {
		if answer.Matches(message) {
			return true
		}
	}
	return false
}

func NewAnswerInput(correctAnswers []AnswerInput) Answer {
	return &AnswerInputs{
		CorrectAnswers: correctAnswers,
	}
}

// answer type

type AnswerType string

func (AnswerType) GormDataType() string { return "varchar(1)" }

var (
	errInvalidAnswerType error = errors.New("invalid answer type")
)

func (answerType *AnswerType) Valid() []error {
	var errors []error
	if *answerType != AnswerInputsAnswerType && *answerType != AnswerOptionsAnswerType {
		errors = append(errors, errInvalidAnswerType)
	}
	return errors
}

func (answerType *AnswerType) RawAnswer(uknown any) Answer {
	bytes, _ := json.Marshal(uknown)
	return answerType.Answer(bytes)
}

func (answerType *AnswerType) Answer(encoded []byte) Answer {
	if len(answerType.Valid()) != 0 {
		return nil
	}
	switch *answerType {
	case AnswerInputsAnswerType:
		var answer AnswerInputs
		json.Unmarshal(encoded, &answer)
		return &answer
	case AnswerOptionsAnswerType:
		var answer AnswerOptions
		json.Unmarshal(encoded, &answer)
		return &answer
	default:
		return nil
	}
}

// question content

type Question string

func (Question) GormDataType() string { return "varchar(64)" }

var (
	errQuestionCannotBeEmpty error = errors.New("question cannot be empty")
	errQuestionToLong        error = errors.New("question cannot exceed 64 characters")
)

func (question *Question) Valid() []error {
	var errors []error
	if *question == "" {
		errors = append(errors, errQuestionCannotBeEmpty)
	}
	if len(*question) > 64 {
		errors = append(errors, errQuestionToLong)
	}
	return errors
}

// question set name

type QuestionSetName string

func (QuestionSetName) GormDataType() string { return "varchar(64)" }

var (
	errQuestionSetNameCannotBeEmpty error = errors.New("question set name cannot be empty")
	errQuestionSetNameToLong        error = errors.New("question set name exceeded 64 characters")
)

func (name *QuestionSetName) Valid() []error {
	var errors []error
	if *name == "" {
		errors = append(errors, errQuestionSetNameCannotBeEmpty)
	}
	if len(*name) > 64 {
		errors = append(errors, errQuestionSetNameToLong)
	}
	return errors
}

// question set description

type QuestionSetDescription string

func (QuestionSetDescription) GormDataType() string { return "varchar(512)" }

var (
	errQuestionSetDescriptionCannotBeEmpty error = errors.New("question set description cannot be empty")
	errQuestionSetDescriptiontoLong        error = errors.New("question set description exceeded 512 characters")
)

func (desc *QuestionSetDescription) Valid() []error {
	var errors []error
	if *desc == "" {
		errors = append(errors, errQuestionSetDescriptionCannotBeEmpty)
	}
	if len(*desc) > 512 {
		errors = append(errors, errQuestionSetDescriptiontoLong)
	}
	return errors
}
