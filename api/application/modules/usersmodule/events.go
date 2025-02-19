package usersmodule

import (
	"quizapi/common"
	"quizapi/modules/eventsmodule"
	"quizapi/modules/modelmodule"
)

var events []eventsmodule.Event = []eventsmodule.Event{
	eventsmodule.NewEvent(DisconnectedEvent{}),
}

type DisconnectedEvent struct {
	Services common.Ioc
	UserId   modelmodule.ModelId
}
