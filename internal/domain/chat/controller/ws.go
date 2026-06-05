package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	handshakeTimeout = 60 * time.Second // Таймаут на хендшейк
	pongWait         = 60 * time.Second // Максимальное время ожидания Pong от клиента
)

// Upgrader - создает сокетовое соединение с удаленным клиентом
func Upgrader(ctx *gin.Context) (*websocket.Conn, error) {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: handshakeTimeout,
	}

	conn, uErr := u.Upgrade(ctx.Writer, ctx.Request, nil)
	if uErr != nil {
		return nil, uErr
	}

	// Устанавливаем настройки чтения
	conn.SetReadLimit(512 * 1024)
	if lErr := conn.SetReadDeadline(time.Now().Add(pongWait)); lErr != nil {
		return nil, lErr
	}

	// При получении Pong, сдвигаем дедлайн
	conn.SetPongHandler(func(string) error {
		if pErr := conn.SetReadDeadline(time.Now().Add(pongWait)); pErr != nil {
			return pErr
		}
		return nil
	})

	return conn, nil
}
