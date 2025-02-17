package modelmodule

import "quizapi/common"

var (
	ErrIdCannotBeEmpty      common.ErrorWithPath = common.NewErrorWithPath("id cannot be empty")
	ErrParallelModification common.HttpError     = common.NewHttpError("repeat. model was modified in the same time by someone else", 409)
)
