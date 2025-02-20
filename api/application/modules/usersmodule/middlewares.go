package usersmodule

import (
	"quizapi/common"
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

	next()
}

func authorizeWsMiddleware(c common.Ioc, _ []byte, next func()) {
	var socketIdStorage common.ServiceStorage[wsmodule.SocketId]
	c.Inject(&socketIdStorage)
	id := socketIdStorage.MustGet()

	var socketRepo UserSocketRepo
	c.Inject(&socketRepo)
	socket := socketRepo.GetBySocket(c, id)
	if socket == nil {
		return
	}

	var userRepo UserRepository
	c.Inject(&userRepo)
	user := userRepo.GetById(c, socket.UserId)
	if user == nil { // user got deleted when in game
		var socketStorage wsmodule.SocketStorage
		c.Inject(&socketStorage)
		for _, us := range socketRepo.GetByUser(c, socket.UserId) {
			socketStorage.Close(us.SocketId)
		}
		return
	}

	session := NewSessionDto(*user)
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)
	sessionStorage.Set(session)

	next()
}
