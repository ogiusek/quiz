package usersmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/wsmodule"

	"github.com/valyala/fasthttp"
)

func authorizeHttpMiddleware(ctx *fasthttp.RequestCtx, c common.Ioc, next func()) {
	authorization := string(ctx.Request.Header.Peek("Authorization"))

	if authorization == "" {
		authorization = string(ctx.QueryArgs().Peek("authorization"))
	}

	if authorization == "" {
		next()
		return
	}

	var session SessionDto
	if err := session.Decode(c, authorization); err != nil {
		ctx.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	}

	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)
	sessionStorage.Set(session)

	var socketIdStorage common.ServiceStorage[wsmodule.SocketId]
	c.Inject(&socketIdStorage)
	socketIdStorage.Set(wsmodule.SocketId(session.UserId))

	next()
}

func authorizeWsMiddleware(c common.Ioc, _ []byte, next func()) {
	var socketIdStorage common.ServiceStorage[wsmodule.SocketId]
	c.Inject(&socketIdStorage)
	id := socketIdStorage.Get()
	if id == nil {
		log.Panic("socket id is null")
	}

	var userRepo UserRepository
	c.Inject(&userRepo)
	user := userRepo.GetById(c, modelmodule.ModelId(*id))
	if user == nil { // user got deleted when in game
		var socketStorage wsmodule.SocketStorage
		c.Inject(&socketStorage)
		socketStorage.Close(*id)
		return
	}

	session := NewSessionDto(*user)
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)
	sessionStorage.Set(session)

	next()
}
