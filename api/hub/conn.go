package main

import (
	"github.com/fasthttp/websocket"
	"github.com/ogiusek/wshub"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(ctx *fasthttp.RequestCtx) bool { return true },
}

type wscon struct {
	conn *websocket.Conn
}

func (conn *wscon) Close() { conn.conn.Close() }
func (conn *wscon) ReadMessage() ([]byte, error) {
	_, p, e := conn.conn.ReadMessage()
	return p, e
}
func (conn *wscon) Send(p []byte)                     { conn.conn.WriteMessage(websocket.TextMessage, p) }
func ToHubConn(conn *websocket.Conn) wshub.SocketConn { return &wscon{conn: conn} }
