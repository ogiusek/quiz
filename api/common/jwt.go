package common

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	errJwtTokenIsEmpty   error = errors.New("missing jwt token")
	errJwtTokenIsInvalid error = errors.New("invalid jwt token")
)

type JwtClaims map[string]any

type JwtConfig struct {
	Secret     interface{}
	SignMethod *jwt.SigningMethodHMAC // jwt.SigningMethodHS256
}

func NewJwtConfig(secret string) JwtConfig {
	return JwtConfig{
		Secret:     []byte(secret), // example secret "npVLKFqcrHOHbwfk84YNmohGNP9vdZtQ"
		SignMethod: jwt.SigningMethodHS256,
	}
}

func (config *JwtConfig) EncodeJwt(now time.Time, claims JwtClaims, expire time.Duration) string {
	claims["exp"] = now.Add(expire).Unix()

	unsignedToken := jwt.NewWithClaims(config.SignMethod, jwt.MapClaims(claims))
	token, err := unsignedToken.SignedString(config.Secret)
	if err != nil {
		log.Panic(err.Error())
	}

	return string(token)
}

func (config *JwtConfig) DecodeJwt(rawToken string) (JwtClaims, error) {
	if rawToken == "" {
		return nil, errJwtTokenIsEmpty
	}

	token, err := jwt.Parse(string(rawToken), func(token *jwt.Token) (interface{}, error) {
		return config.Secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errJwtTokenIsInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Panic("this should always be claims but it isn't")
	}

	return JwtClaims(claims), nil
}
