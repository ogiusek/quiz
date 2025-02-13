package usersmodule

import (
	"errors"
	"fmt"
	"quizapi/common"
	"quizapi/modules/timemodule"
	"time"
)

// config

type UserConfig struct {
	SessionTokenExpirationTime timemodule.Duration `json:"session_token_expiration_time"`
	RefreshTokenExpirationTime timemodule.Duration `json:"refresh_token_expiration_time"`
}

func (o *UserConfig) Valid() []error {
	var errs []error
	if o.SessionTokenExpirationTime.Duration().Seconds() < 5 {
		errs = append(errs, common.ErrPath(errors.New("session token expiration time cannot be shorter than 5 seconds")).Property("session_token_expiration_time"))
	}
	if o.RefreshTokenExpirationTime.Duration().Seconds() < o.SessionTokenExpirationTime.Duration().Seconds() {
		errs = append(errs, common.ErrPath(errors.New("refresh token expiration time cannot be shorter than session token expiration time")).Property("refresh_token_expiration_time"))
	}
	if len(errs) != 0 {
		errs = append(errs, fmt.Errorf("remember example time date format is %s", (time.Hour*(0)+time.Minute*0+time.Second*60).String()))
	}
	return errs
}

const ( // tokens types
	sessionTokenType string = "session_token"
	refreshTokenType string = "refresh_token"
)
