package common

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/valyala/fasthttp"
)

// requires:
// - http 			// can run away
// - validable 	// can run away
// - ioc 				// will be always restricted to using this

type HttpMiddleware interface {
	Handle(ctx *fasthttp.RequestCtx, c Ioc, next func())
}

type httpMiddlewareImpl struct {
	handle func(ctx *fasthttp.RequestCtx, c Ioc, next func())
}

func (m *httpMiddlewareImpl) Handle(ctx *fasthttp.RequestCtx, c Ioc, next func()) {
	m.handle(ctx, c, next)
}

func NewHttpMiddleware(handle func(ctx *fasthttp.RequestCtx, c Ioc, next func())) HttpMiddleware {
	return &httpMiddlewareImpl{
		handle: handle,
	}
}

func RunHttpMiddlewares(middlewares []HttpMiddleware, endpoint func(ctx *fasthttp.RequestCtx, c Ioc), ctx *fasthttp.RequestCtx, c Ioc) {
	var next func()
	next = func() {
		if len(middlewares) == 0 {
			endpoint(ctx, c)
			return
		}
		middleware := middlewares[0]
		middlewares = middlewares[1:]
		middleware.Handle(ctx, c, next)
	}

	next()
}

//

type Response any

type Endpoint interface {
	Handle(c Ioc) error
}

func FromBody(ctx *fasthttp.RequestCtx, modelPtr any) error {
	return json.Unmarshal(ctx.Request.Body(), modelPtr)
}

func FromQuery(ctx *fasthttp.RequestCtx, modelPtr any) error {
	args := ctx.QueryArgs().Peek("args")
	if string(args) == "" {
		return nil
	}
	err := json.Unmarshal(args, modelPtr)
	return err
}

// argsPtr have to be a pointer to empty args model
func HttpEndpoint(c IocScope, modelBinder func(ctx *fasthttp.RequestCtx, modelPtr any) error, argsPtr any) func(ctx *fasthttp.RequestCtx) {
	// endpoint method arguments and result shouldn't change.
	// modifications here and not necessary now
	// but if you're going to work here i reccomended some changes

	var middlewaresGroup ServiceGroup[HttpMiddleware]
	c.Scope().Inject(&middlewaresGroup)
	var middlewares = middlewaresGroup.GetAll()

	argsType := reflect.TypeOf(argsPtr).Elem()
	if _, ok := argsPtr.(Endpoint); !ok {
		log.Panic("args which do not match 'Endpoint' interface cannot be mapped to endpoint")
	}

	endpoint := func(ctx *fasthttp.RequestCtx, c Ioc) {
		args := reflect.New(argsType).Interface()
		if err := modelBinder(ctx, &args); err != nil {
			ctx.SetBody([]byte(err.Error()))
			ctx.SetStatusCode(415)
			return
		}

		// should be extracted
		// TODO
		// create ServiceGroup[ModelValidation]
		validable, ok := args.(Validable)
		if ok {
			if errs := validable.Valid(); len(errs) != 0 {
				var messages []string
				for _, err := range errs {
					messages = append(messages, ErrPath(err).Error())
				}
				ctx.SetStatusCode(400)
				ctx.SetBody([]byte(strings.Join(messages, "\n")))
				return
			}
		}

		err := args.(Endpoint).Handle(c)

		if err == nil {
			// TODO
			// should be middleware
			var responseStorage ServiceStorage[Response]
			c.Inject(&responseStorage)
			response := responseStorage.Get()

			if response == nil {
				ctx.SetStatusCode(204)
				return
			}

			ctx.Response.Header.Add("Content-Type", "application/json")
			res, err := json.Marshal(*response)
			if err != nil {
				log.Panic(err.Error())
			}
			ctx.SetBody(res)

			return
		}

		var errStorage ServiceStorage[error]
		c.Inject(&errStorage)
		errStorage.Set(err)

		// should be also extracted
		// TODO
		// create ServiceGroup[ErrHandler] where "type ErrHandler func(ctx, error) bool" (bool means ok and stops iteration)
		if httpError, ok := err.(HttpError); ok {
			ctx.SetStatusCode(httpError.StatusCode())
			ctx.SetBody([]byte(httpError.Error()))
			return
		}

		// this can be left here if no ErrHandler were found
		log.Print(err.Error())
		ctx.SetStatusCode(500)
	}

	return func(ctx *fasthttp.RequestCtx) {
		scope := c.Scope()
		RunHttpMiddlewares(middlewares, endpoint, ctx, scope)
	}
}
