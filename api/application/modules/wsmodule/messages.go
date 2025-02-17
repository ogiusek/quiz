package wsmodule

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type MessageTopic string
type MessagePayload []byte

func emptyPayload() []byte {
	p, _ := json.Marshal(struct{}{})
	return p
}

var EmptyPayload []byte = emptyPayload()

type Message struct {
	Topic   MessageTopic `json:"topic"`
	Payload any          `json:"payload"`
}

func NewMessage(topic string, payload any) Message {
	return Message{Topic: MessageTopic(topic), Payload: payload}
}

type SocketsMessager interface {
	// panics when there is already existing listener
	Listen(MessageTopic, func(SocketId, SocketConn, MessagePayload))

	GetListenersTopics() []MessageTopic

	// passes Sockets error
	Send(SocketId, Message) error
}

type socketsMessagerImpl struct {
	sockets                SocketStorage
	listeners              map[MessageTopic]func(SocketId, SocketConn, MessagePayload)
	missingPayloadListener func(SocketId)
	cannotDecodeListener   func(SocketId, []byte, error)
	invalidTopicListener   func(SocketId, Message)
}

func (messager *socketsMessagerImpl) listen() {
	messager.sockets.OnMessage(func(id SocketId, conn SocketConn, bytes []byte) {
		var decodedMessage struct {
			Topic   MessageTopic `json:"topic"`
			Payload any          `json:"payload"`
		}

		if err := json.Unmarshal(bytes, &decodedMessage); err != nil {
			messager.cannotDecodeListener(id, bytes, err)
			return
		}

		if decodedMessage.Payload == nil {
			messager.missingPayloadListener(id)
			return
		}

		payload, _ := json.Marshal(decodedMessage.Payload)
		message := Message{
			Topic:   decodedMessage.Topic,
			Payload: payload,
		}

		listener, ok := messager.listeners[decodedMessage.Topic]
		if !ok {
			messager.invalidTopicListener(id, message)
			return
		}

		listener(id, conn, payload)
	})
}

var (
	errInvalidMessage error = errors.New("invalid message")
	errMissingPayload error = errors.New("missing payload")
)

func NewSocketsMessager(sockets SocketStorage) SocketsMessager {
	var messager *socketsMessagerImpl
	messager = &socketsMessagerImpl{
		sockets:   sockets,
		listeners: make(map[MessageTopic]func(SocketId, SocketConn, MessagePayload)),
		cannotDecodeListener: func(socket SocketId, message []byte, err error) {
			sockets.SendMessage(socket, []byte(errInvalidMessage.Error()))
		},
		invalidTopicListener: func(socket SocketId, message Message) {
			response := []byte(fmt.Sprintf("'%s' topic listener does not exist", message.Topic))
			messager.sockets.SendMessage(socket, response)
		},
		missingPayloadListener: func(id SocketId) {
			messager.sockets.SendMessage(id, []byte(errMissingPayload.Error()))
		},
	}
	messager.listen()
	return messager
}

func (m *socketsMessagerImpl) Listen(topic MessageTopic, handler func(SocketId, SocketConn, MessagePayload)) {
	if _, exists := m.listeners[topic]; exists {
		log.Panicf("cannot add '%s' listener because it already exists", topic)
	}
	m.listeners[topic] = handler
}

func (m *socketsMessagerImpl) GetListenersTopics() []MessageTopic {
	var keys []MessageTopic
	for key := range m.listeners {
		keys = append(keys, key)
	}
	return keys
}

func (m *socketsMessagerImpl) Send(id SocketId, message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return m.sockets.SendMessage(id, bytes)
}
