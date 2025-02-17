package wsmodule

import (
	"encoding/json"
	"log"
	"quizapi/common"
	"reflect"
	"strings"

	"github.com/fasthttp/websocket"
)

type Middleware interface {
	Invoke(c common.Ioc, message []byte, next func())
}

type middleware func(c common.Ioc, message []byte, next func())

func (m *middleware) Invoke(c common.Ioc, message []byte, next func()) {
	(*m)(c, message, next)
}

func RunMiddlewares(middlewares []Middleware, c common.Ioc, message []byte, endpoint func(c common.Ioc, message []byte)) {
	var next func()
	next = func() {
		if len(middlewares) == 0 {
			endpoint(c, message)
			return
		}
		middleware := middlewares[0]
		middlewares = middlewares[1:]
		middleware.Invoke(c, message, next)
	}
	next()
}

func NewMiddleware(invoke func(c common.Ioc, message []byte, next func())) Middleware {
	m := middleware(invoke)
	return &m
}

// argsPtr have to be a pointer to empty args model
func WsEndpoint(c common.IocScope, argsPtr any) func(SocketId, SocketConn, MessagePayload) {
	argsType := reflect.TypeOf(argsPtr).Elem()
	if _, ok := argsPtr.(common.Endpoint); !ok {
		var logger log.Logger
		c.Scope().Inject(&logger)
		logger.Panic("args which do not match 'Endpoint' interface cannot be mapped to endpoint")
	}

	var sockets SocketStorage
	var socketMessager SocketsMessager
	c.Scope().Inject(&socketMessager)
	c.Scope().Inject(&sockets)

	endpoint := func(c common.Ioc, bytes []byte) {
		var connStorage common.ServiceStorage[SocketConn]
		c.Inject(&connStorage)
		conn := connStorage.MustGet()

		args := reflect.New(argsType).Interface()

		if err := json.Unmarshal(bytes, &args); err != nil {
			conn.conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		validable, ok := args.(common.Validable)
		if ok {
			if errs := validable.Valid(); len(errs) != 0 {
				var messages []string
				for _, err := range errs {
					messages = append(messages, common.ErrPath(err).Property("payload").Error())
				}

				response := strings.Join(messages, "\n")
				conn.conn.WriteMessage(websocket.TextMessage, []byte(response))
				return
			}
		}

		err := args.(common.Endpoint).Handle(c)

		if err == nil {
			var responseStorage common.ServiceStorage[common.Response]
			c.Inject(&responseStorage)
			response := responseStorage.Get()

			if response == nil {
				return
			}

			res, err := json.Marshal(*response)
			if err != nil {
				var logger log.Logger
				c.Inject(&logger)
				logger.Panicf("failure when encoding response %v, %s", *response, err.Error())
			}

			conn.conn.WriteMessage(websocket.TextMessage, res)
			return
		}

		var errStorage common.ServiceStorage[error]
		c.Inject(&errStorage)
		errStorage.Set(err)

		if httpError, ok := err.(common.HttpError); ok {
			conn.conn.WriteMessage(websocket.TextMessage, []byte(httpError.Error()))
			return
		}

		var logger log.Logger
		c.Inject(&logger)
		logger.Print(err.Error())
		conn.conn.WriteMessage(websocket.TextMessage, []byte("server error"))
	}

	var middlewaresGroup common.ServiceGroup[Middleware]
	c.Scope().Inject(&middlewaresGroup)
	middlewares := middlewaresGroup.GetAll()

	return func(id SocketId, conn SocketConn, m MessagePayload) {
		scope := c.Scope()
		var socketIdStorage common.ServiceStorage[SocketId]
		scope.Inject(&socketIdStorage)
		socketIdStorage.Set(id)

		var socketConnStorage common.ServiceStorage[SocketConn]
		scope.Inject(&socketConnStorage)
		socketConnStorage.Set(conn)

		RunMiddlewares(middlewares, scope, m, endpoint)
	}
}
