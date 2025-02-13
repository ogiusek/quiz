package usersmodule

import "time"

type UserConfig struct {
	SessionTokenExpirationTime time.Duration
	RefreshTokenExpirationTime time.Duration
}

const ( // tokens types
	sessionTokenType string = "session_token"
	refreshTokenType string = "refresh_token"
)
