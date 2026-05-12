package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	handshakeTimeout = 60 * time.Second
)

// Upgrader - создает сокетовое соединение с удаленным клиентом
func Upgrader(ctx *gin.Context) (*websocket.Conn, error) {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		HandshakeTimeout: handshakeTimeout,
	}

	conn, err := u.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(8192)

	return conn, nil
}
