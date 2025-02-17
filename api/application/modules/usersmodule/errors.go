package usersmodule

import (
	"quizapi/common"
)

var (
	ErrUnauthorized common.HttpError = common.NewHttpError("missing valid authorization", 401)
	ErrUserNotFound common.HttpError = common.NewHttpError("user not found", 404)
)

var (
	errMissingSessionToken    common.HttpError = common.NewHttpError("missing session token", 400)
	errMissingRefreshToken    common.HttpError = common.NewHttpError("missing refresh token", 400)
	errNotASessionToken       common.HttpError = common.NewHttpError("this is not a session token", 400)
	errNotARefreshToken       common.HttpError = common.NewHttpError("this is not a refresh token", 400)
	errFabricatedSessionToken common.HttpError = common.NewHttpError("fabricated session token", 400)

	errUserConflict           common.HttpError = common.NewHttpError("this name is already taken", 409) // only name is unique
	errUserInvalidCredentials common.HttpError = common.NewHttpError("user with this credentials do not exists", 400)
)
