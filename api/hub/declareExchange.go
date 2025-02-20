package main

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	ExchangeWsStarted             string = "ws_started"
	ExchangeWsConnectRequest      string = "ws_connect_request"
	ExchangeWsConnectConfirmation string = "ws_connect_confirmation"
	ExchangeWsReceived            string = "ws_received"
	ExchangeWsRespond             string = "ws_respond"
	ExchangeWsClose               string = "ws_close"
	ExchangeWsClosed              string = "ws_closed"
)

const (
	QueueWsStarted             string = "ws_main_started"
	QueueWsConnectRequest      string = "ws_main_connect_request"
	QueueWsConnectConfirmation string = "ws_main_connect_confirmation"
	QueueWsReceived            string = "ws_main_received"
	QueueWsRespond             string = "ws_main_respond"
	QueueWsClose               string = "ws_main_close"
	QueueWsClosed              string = "ws_main_closed"
)

var binds map[string]string = map[string]string{
	ExchangeWsStarted:             QueueWsStarted,
	ExchangeWsConnectRequest:      QueueWsConnectRequest,
	ExchangeWsConnectConfirmation: QueueWsConnectConfirmation,
	ExchangeWsReceived:            QueueWsReceived,
	ExchangeWsRespond:             QueueWsRespond,
	ExchangeWsClose:               QueueWsClose,
	ExchangeWsClosed:              QueueWsClosed,
}

func DeclareExchanges(ch *amqp.Channel) {
	for exchange, queue := range binds {
		if err := ch.ExchangeDeclare(exchange, amqp.ExchangeDirect, false, false, false, false, nil); err != nil {
			log.Panic(err.Error())
		}

		if _, err := ch.QueueDeclare(queue, false, false, false, false, nil); err != nil {
			log.Panic(err.Error())
		}

		if err := ch.QueueBind(queue, "", exchange, false, nil); err != nil {
			log.Panic(err.Error())
		}
	}
}
