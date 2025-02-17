package eventsmodule

import "reflect"

type Event struct {
	Topic   string
	Payload any
}

func NewEvent[T any](payload T) Event {
	return Event{
		Topic:   reflect.TypeOf(payload).String(),
		Payload: payload,
	}
}
