package questionsmodule

import (
	"quizapi/common"
)

var (
	errQuestionSetConflict  error = common.NewHttpError("question set with this name already exists", 409)
	errQuestionConflict     error = common.NewHttpError("question conflict", 409)
	ErrQuestionSetNotFound  error = common.NewHttpError("cannot find question set", 404)
	errQuestionNotFound     error = common.NewHttpError("cannot find question", 404)
	errQuestionSetForbidden error = common.NewHttpError("you are not the owner of this question set", 403)
)
