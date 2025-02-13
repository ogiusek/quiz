package wsmodule

import (
	"errors"
	"sync"

	"github.com/fasthttp/websocket"
)

type SocketId string
type SocketConn struct {
	conn *websocket.Conn
}

var (
	ErrSocketNotFound error = errors.New("this socket do not exists")
	ErrSocketConflict error = errors.New("this socket already exists")
)

type socketConnect interface {
	IsTaken(socket SocketId) bool
	Connect(socket SocketId, conn *websocket.Conn) error
}

type SocketStorage interface {
	SendMessage(socket SocketId, message []byte) error
	Close(socket SocketId) error

	OnConnect(listener func(id SocketId, conn SocketConn))
	OnMessage(listener func(id SocketId, conn SocketConn, message []byte))
	OnClose(listener func(id SocketId, conn SocketConn))
}

type socketStorageImpl struct {
	mux              sync.Mutex
	sockets          map[SocketId][]*websocket.Conn
	connectListeners []func(SocketId, SocketConn)
	messageListeners []func(SocketId, SocketConn, []byte)
	closeListeners   []func(SocketId, SocketConn)
}

func NewSockets() SocketStorage {
	return &socketStorageImpl{
		sockets: map[SocketId][]*websocket.Conn{},
	}
}

func (impl *socketStorageImpl) getConn(id SocketId) ([]*websocket.Conn, bool) {
	conn, ok := impl.sockets[id]
	return conn, ok
}

func (impl *socketStorageImpl) run(id SocketId, conn *websocket.Conn) {
	defer conn.Close()
	impl.mux.Lock()
	impl.sockets[id] = append(impl.sockets[id], conn)
	impl.mux.Unlock()

	for _, listener := range impl.connectListeners {
		listener(id, SocketConn{conn: conn})
	}

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		for _, listener := range impl.messageListeners {
			listener(id, SocketConn{conn: conn}, bytes)
		}
	}

	impl.mux.Lock()
	sockets := []*websocket.Conn{}
	for _, existingConn := range impl.sockets[id] {
		if existingConn != conn {
			sockets = append(sockets, existingConn)
		}
	}
	impl.sockets[id] = sockets
	impl.mux.Unlock()

	if len(impl.sockets[id]) == 0 {
		for _, listener := range impl.closeListeners {
			listener(id, SocketConn{conn: conn})
		}
	}
}

func (impl *socketStorageImpl) IsTaken(id SocketId) bool {
	_, exists := impl.sockets[id]
	return exists
}

func (sockets *socketStorageImpl) Connect(id SocketId, conn *websocket.Conn) error {
	// if _, ok := sockets.sockets[id]; ok {
	// 	return ErrSocketConflict
	// }
	// sockets.sockets[id]
	sockets.run(id, conn)
	return nil
}

func (sockets *socketStorageImpl) SendMessage(id SocketId, message []byte) error {
	conns, ok := sockets.getConn(id)
	if !ok {
		return ErrSocketNotFound
	}
	for _, conn := range conns {
		conn.WriteMessage(websocket.TextMessage, message)
	}
	return nil
}

func (sockets *socketStorageImpl) Close(id SocketId) error {
	conns, ok := sockets.getConn(id)
	if !ok {
		return ErrSocketNotFound
	}
	for _, conn := range conns {
		conn.Close()
	}
	return nil
}

func (sockets *socketStorageImpl) OnConnect(listener func(SocketId, SocketConn)) {
	sockets.connectListeners = append(sockets.connectListeners, listener)
}

func (sockets *socketStorageImpl) OnMessage(listener func(SocketId, SocketConn, []byte)) {
	sockets.messageListeners = append(sockets.messageListeners, listener)
}

func (sockets *socketStorageImpl) OnClose(listener func(SocketId, SocketConn)) {
	sockets.closeListeners = append(sockets.closeListeners, listener)
}
