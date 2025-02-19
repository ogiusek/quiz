package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fasthttp/websocket"
	"github.com/ogiusek/wshub"
	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

func main() {
	config := GetConfig()

	conn, err := amqp.Dial(config.RabbitMqUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	DeclareExchanges(ch)

	startedSender := wshub.NewBrokerSender(func(m wshub.Started) {
		body, _ := json.Marshal(m)
		ch.Publish(ExchangeWsStarted, "", false, false, amqp.Publishing{Body: body})
	})

	connectRequest := wshub.NewBrokerSender(func(m wshub.ConnectRequest) {
		body, _ := json.Marshal(m)
		ch.Publish(ExchangeWsConnectRequest, "", false, false, amqp.Publishing{Body: body})
	})

	connectConfirmationListener, connectConfirmation := wshub.NewBrokerListener[wshub.ConnectConfirmation]()

	received := wshub.NewBrokerSender(func(m wshub.SocketMessage) {
		body, _ := json.Marshal(m)
		ch.Publish(ExchangeWsReceived, "", false, false, amqp.Publishing{Body: body})
	})

	respondListener, respond := wshub.NewBrokerListener[wshub.SocketMessage]()

	closed := wshub.NewBrokerSender(func(m wshub.Close) {
		body, _ := json.Marshal(m)
		ch.Publish(ExchangeWsClosed, "", false, false, amqp.Publishing{Body: body})
	})

	closeListener, close := wshub.NewBrokerListener[wshub.Close]()

	if messages, err := ch.Consume(QueueWsConnectConfirmation, "", false, false, false, false, nil); err != nil {
		log.Panic()
	} else {
		go func() {
			for {
				message := <-messages
				message.Ack(false)

				var connectConfirmationMessage wshub.ConnectConfirmation
				json.Unmarshal(message.Body, &connectConfirmationMessage)
				connectConfirmation(connectConfirmationMessage)
			}
		}()
	}

	if messages, err := ch.Consume(QueueWsRespond, "", false, false, false, false, nil); err != nil {
		log.Panic()
	} else {
		go func() {
			for {
				message := <-messages
				message.Ack(false)

				var socketMessage wshub.SocketMessage
				json.Unmarshal(message.Body, &socketMessage)
				respond(socketMessage)
			}
		}()
	}

	if messages, err := ch.Consume(QueueWsClose, "", false, false, false, false, nil); err != nil {
		log.Panic()
	} else {
		go func() {
			for {
				message := <-messages
				message.Ack(false)

				var closeMessage wshub.Close
				json.Unmarshal(message.Body, &closeMessage)
				close(closeMessage)
			}
		}()
	}

	broker := wshub.NewBroker(
		startedSender,
		connectRequest,
		connectConfirmationListener,
		received,
		respondListener,
		closeListener,
		closed,
	)

	hub := wshub.NewWsHub(broker)

	log.Printf("starting server on :%d", config.Port)
	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", config.Port), func(ctx *fasthttp.RequestCtx) {
		// CORS
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type")

		// handle request
		auth := ctx.Request.URI().QueryArgs().Peek("authorization")

		// waiter is required because of the order in which things has to happen
		// 1. hub.Connect argument func has to be triggered
		// 2. upgrader.Upgrade has to be called
		// 3. endpoint has to finish to trigger upgrader.Upgrade argument func
		// 4. hub.Connect argument has to finish
		// 5. hub.Connect ends after websocket closes
		// 6. upgrader.Upgrade has to finish
		waiter := make(chan bool, 1)

		go func() {
			// this works until socket stops
			hub.Connect(func() wshub.SocketConn {
				connResolver := make(chan wshub.SocketConn, 1)
				// when this stops websocket closes
				upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
					connResolver <- ToHubConn(conn)
					<-waiter
				})
				waiter <- true
				conn := <-connResolver
				return conn
			}, auth)
			waiter <- true
		}()

		<-waiter

		// when this ends upgrader trigger
	}); err != nil {
		log.Panic(err.Error())
	}
}
