package wsmodule

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"quizapi/common"

	"github.com/fasthttp/websocket"
	"github.com/ogiusek/wshub"
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

//

type SocketId wshub.Id

// Implement the driver.Valuer interface.
func (v SocketId) Value() (driver.Value, error) {
	id := wshub.Id(v)
	return id.String(), nil
}

// Implement the driver.Scanner interface.
func (v *SocketId) Scan(value interface{}) error {
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string value")
	}

	*v = SocketId(wshub.SocketIdFrom(stringValue))

	return nil
}

var (
	ErrSocketNotFound error = errors.New("this socket do not exists")
	ErrSocketConflict error = errors.New("this socket already exists")
)

type socketConnect interface {
	Listen()
}

type SocketStorage interface {
	SendMessage(socket SocketId, message []byte)
	Close(socket SocketId)

	OnStart(listener func())
	OnConnect(listener func(id SocketId))
	OnMessage(listener func(id SocketId, message []byte))
	OnClose(listener func(id SocketId))
}

type socketInterface interface {
	socketConnect
	SocketStorage
}

type socketStorageImpl struct {
	channel          *amqp.Channel
	connect          func(socketId SocketId, meta []byte) (canConnect bool)
	sockets          map[SocketId][]*websocket.Conn
	startListeners   []func()
	connectListeners []func(SocketId)
	messageListeners []func(SocketId, []byte)
	closeListeners   []func(SocketId)
}

func NewSocketStorage(c common.Ioc, connectHandler func(socketId SocketId, meta []byte) (canConnect bool)) socketInterface {
	var ch *amqp.Channel
	c.Inject(&ch)
	storage := &socketStorageImpl{
		channel: ch,
		connect: connectHandler,
		sockets: map[SocketId][]*websocket.Conn{},
	}

	return storage
}

func (storage *socketStorageImpl) Listen() {
	go func() {
		messages, err := storage.channel.Consume(QueueWsStarted, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err.Error())
		}
		for {
			message := <-messages
			for _, listener := range storage.startListeners {
				listener()
			}
			message.Ack(true)
		}
	}()
	go func() {
		messages, err := storage.channel.Consume(QueueWsConnectRequest, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err.Error())
		}
		for {
			message := <-messages
			var body wshub.ConnectRequest
			json.Unmarshal(message.Body, &body)
			canConnect := storage.connect(SocketId(body.SocketId), body.Payload)
			for _, listener := range storage.connectListeners {
				listener(SocketId(body.SocketId))
			}
			connectResponse := wshub.NewConnectConfirmation(body.SocketId, canConnect)
			resBody, _ := json.Marshal(connectResponse)
			storage.channel.Publish(ExchangeWsConnectConfirmation, "", false, false, amqp.Publishing{Body: resBody})
			message.Ack(false)
		}
	}()
	go func() {
		messages, err := storage.channel.Consume(QueueWsReceived, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err.Error())
		}
		for {
			message := <-messages
			var body wshub.SocketMessage
			json.Unmarshal(message.Body, &body)
			for _, listener := range storage.messageListeners {
				listener(SocketId(body.SocketId), body.Payload)
			}
			message.Ack(false)
		}
	}()
	go func() {
		messages, err := storage.channel.Consume(QueueWsClosed, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err.Error())
		}
		for {
			message := <-messages
			var body wshub.Close
			json.Unmarshal(message.Body, &body)
			for _, listener := range storage.closeListeners {
				listener(SocketId(body.SocketId))
			}
			message.Ack(false)
		}
	}()
}

func (sockets *socketStorageImpl) SendMessage(id SocketId, bytes []byte) {
	message := wshub.NewSocketMessage(wshub.Id(id), bytes)
	body, _ := json.Marshal(message)
	sockets.channel.Publish(ExchangeWsRespond, "", false, false, amqp.Publishing{Body: body})
}

func (sockets *socketStorageImpl) Close(id SocketId) {
	message := wshub.NewClose(wshub.Id(id))
	body, _ := json.Marshal(message)
	sockets.channel.Publish(ExchangeWsClose, "", false, false, amqp.Publishing{Body: body})
}

func (sockets *socketStorageImpl) OnStart(listenr func()) {
	sockets.startListeners = append(sockets.startListeners, listenr)
}

func (sockets *socketStorageImpl) OnConnect(listener func(SocketId)) {
	sockets.connectListeners = append(sockets.connectListeners, listener)
}

func (sockets *socketStorageImpl) OnMessage(listener func(SocketId, []byte)) {
	sockets.messageListeners = append(sockets.messageListeners, listener)
}

func (sockets *socketStorageImpl) OnClose(listener func(SocketId)) {
	sockets.closeListeners = append(sockets.closeListeners, listener)
}
