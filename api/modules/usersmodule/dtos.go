package usersmodule

import (
	"quizapi/common"
	"quizapi/modules/filesmodule"
	"quizapi/modules/modelmodule"
	"quizapi/modules/timemodule"
)

// tokens dto
// session dto
// refresh dto
// user model dto

// tokens dto

type TokensDto struct {
	SessionToken string `json:"session_token"`
	RefreshToken string `json:"refresh_token"`
}

func (tokens *TokensDto) Valid() []error {
	var errors []error
	if tokens.SessionToken == "" {
		errors = append(errors, common.ErrPath(errMissingSessionToken).Property("session_token"))
	}
	if tokens.RefreshToken == "" {
		errors = append(errors, common.ErrPath(errMissingRefreshToken).Property("refresh_token"))
	}
	return errors
}

// session

type SessionDto struct {
	UserId    modelmodule.ModelId `json:"id"`
	UserName  UserName            `json:"name"`
	UserImage filesmodule.ImageId `json:"image"`
}

func NewSessionDto(user UserModel) SessionDto {
	return SessionDto{
		UserId:    user.Id,
		UserName:  user.Name,
		UserImage: user.Image,
	}
}

func (session *SessionDto) Decode(c common.Ioc, sessionToken string) error {
	var config common.JwtConfig
	c.Inject(&config)
	claims, err := config.DecodeJwt(sessionToken)
	if err != nil {
		return err
	}

	tokenType, _ := claims["token_type"].(string)
	if tokenType != sessionTokenType {
		return errNotASessionToken
	}

	session.UserId = modelmodule.ModelId(claims["user_id"].(string))
	session.UserName = UserName(claims["user_name"].(string))
	session.UserImage = filesmodule.NewImageId(filesmodule.FileId(claims["user_image"].(string)))
	return nil
}

func (session *SessionDto) Encode(c common.Ioc) string {
	var jwtConfig common.JwtConfig
	c.Inject(&jwtConfig)

	var userConfig UserConfig
	c.Inject(&userConfig)

	claims := common.JwtClaims{}
	claims["token_type"] = sessionTokenType
	claims["user_id"] = string(session.UserId)
	claims["user_name"] = string(session.UserName)
	claims["user_image"] = string(session.UserImage)
	var clock timemodule.Clock
	c.Inject(&clock)
	return jwtConfig.EncodeJwt(clock.Now(), claims, userConfig.SessionTokenExpirationTime)
}

// refresh dto

type RefreshDto struct {
	SessionHash string              `json:"session_hash"`
	UserId      modelmodule.ModelId `json:"user_id"`
}

func NewRefreshDto(c common.Ioc, sessionToken string, user UserModel) RefreshDto {
	var hasher common.Hasher
	c.Inject(&hasher)
	return RefreshDto{
		SessionHash: hasher.Hash([]byte(sessionToken)),
		UserId:      user.Id,
	}
}

func (refreshDto *RefreshDto) Matches(c common.Ioc, sessionToken string) bool {
	var hasher common.Hasher
	c.Inject(&hasher)
	return refreshDto.SessionHash == hasher.Hash([]byte(sessionToken))
}

func (refreshDto *RefreshDto) Decode(c common.Ioc, refreshToken string) error {
	var config common.JwtConfig
	c.Inject(&config)
	claims, err := config.DecodeJwt(refreshToken)
	if err != nil {
		return err
	}
	tokenType, _ := claims["token_type"].(string)
	if tokenType != refreshTokenType {
		return errNotARefreshToken
	}
	refreshDto.UserId = modelmodule.ModelId(claims["user_id"].(string))
	refreshDto.SessionHash, _ = claims["session_hash"].(string)
	return nil
}

func (refreshDto *RefreshDto) Encode(c common.Ioc) string {
	var jwtConfig common.JwtConfig
	c.Inject(&jwtConfig)

	var userConfig UserConfig
	c.Inject(&userConfig)

	claims := common.JwtClaims{}
	claims["token_type"] = refreshTokenType
	claims["user_id"] = string(refreshDto.UserId)
	claims["session_hash"] = string(refreshDto.SessionHash)
	var clock timemodule.Clock
	c.Inject(&clock)
	return jwtConfig.EncodeJwt(clock.Now(), claims, userConfig.RefreshTokenExpirationTime)
}

// user model dto

type UserDto struct {
	modelmodule.ModelDto
	UserName  UserName            `json:"name"`
	UserImage filesmodule.ImageId `json:"image"`
}

func (user *UserModel) Dto() UserDto {
	return UserDto{
		ModelDto:  user.Model.Dto(),
		UserName:  user.Name,
		UserImage: user.Image,
	}
}
