package wsmodule

import (
	"fmt"
	"quizapi/common"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

// connect
// example args

// connect

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(ctx *fasthttp.RequestCtx) bool { return true },
}

var errAlreadyConnectedError common.HttpError = common.NewHttpError("you are connected on other account", 409)

func Connect(ctx *fasthttp.RequestCtx, c common.Ioc) {
	var connector socketConnect
	c.Inject(&connector)

	var socketStorage common.ServiceStorage[SocketId]
	c.Inject(&socketStorage)
	socket := socketStorage.Get()

	if socket == nil {
		ctx.SetStatusCode(400)
		ctx.SetBodyString("does not meet server expectations. follow documentation to connect")
		return
	}

	upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		err := connector.Connect(SocketId(*socket), conn)
		if err != ErrSocketConflict {
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte(errAlreadyConnectedError.Error()))
		conn.Close()
	})
}

// example args

type Args struct {
	Message string `json:"message"`
}

func (args *Args) Valid() []error {
	var errors []error
	if args.Message == "" {
		errors = append(errors, common.NewErrorWithPath("message cannot be empty").Property("message"))
	}
	return errors
}

func (args *Args) Handle(c common.Ioc) error {
	var respStorage common.ServiceStorage[common.Response]
	c.Inject(&respStorage)
	respStorage.Set(Message{
		Topic:   "kys",
		Payload: fmt.Sprintf("yea i received %s and what ?", args.Message),
	})

	return nil
}
