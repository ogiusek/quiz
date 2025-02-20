package wsmodule

import (
	"fmt"
	"quizapi/common"
)

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
