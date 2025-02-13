package matchmodule

import "quizapi/common"

var (
	errHttpAlreadyInMatch         common.HttpError = common.NewHttpError("you are already in match", 409)
	errHttpPlayerNotFound         common.HttpError = common.NewHttpError("you are not in match", 400)
	errHttpMatchDoNotExists       common.HttpError = common.NewHttpError("match with this id do not exists", 404)
	errHttpMatchCourseDoNotExists common.HttpError = common.NewHttpError("match hasn't started yet", 400)
)
