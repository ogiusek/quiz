package eventsmodule

import (
	"log"
)

type EventManager interface {
	// panics if topic is not reserved
	Listen(topic string, handler func(any))
	// panics if topic is not reserved
	Dispach(event Event)
	// reserves topic
	Reserve(topic string)
}

type eventManager struct {
	listeners map[string][]func(any)
}

func (manager *eventManager) Listen(topic string, handler func(any)) {
	if _, ok := manager.listeners[topic]; !ok {
		log.Panic("topic is not reserved")
	}

	manager.listeners[topic] = append(manager.listeners[topic], handler)
}

func (manager *eventManager) Dispach(event Event) {
	if _, ok := manager.listeners[event.Topic]; !ok {
		log.Panic("topic is not reserved")
	}

	for _, listener := range manager.listeners[event.Topic] {
		listener(event.Payload)
	}
}

func (manager *eventManager) Reserve(topic string) {
	if _, ok := manager.listeners[topic]; ok {
		log.Panicf("topic '%s' is already reserved", topic)
	}

	manager.listeners[topic] = []func(any){}
}

func NewEventManger() EventManager {
	return &eventManager{
		listeners: map[string][]func(any){},
	}
}
