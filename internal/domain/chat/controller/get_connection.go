package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"nofelet/decorator"
	"nofelet/internal/hub"
	"nofelet/middleware/metrics"
)

var errChatRoomNotFound = errors.New("chat room not found")

// GetChat - /chat создает чат
func (c *Controller) GetChat(ctx *gin.Context) {
	// Проверяем uuid комнаты
	uuid := ctx.Param("uuid")
	if len(uuid) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errChatRoomNotFound)
		return
	}

	uConn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.log.Error("socket creation", "err", sErr)
	}
	metrics.WebSocketConnectionsActive.Add(ctx, 1)
	defer metrics.WebSocketConnectionsActive.Add(ctx, -1)

	// Создаем комнату
	chat := hub.NewChatRoom()
	chat.Init(uuid)
	metrics.RoomsActive.Add(ctx, 1)

	// Создаем декоратор коннекшена
	dConn := decorator.New(uConn, c.log, c.cfg)
	defer func() {
		if err := dConn.Close(); err != nil {
			c.log.Error("socket close error", "err", err)
		}
	}()

	// Юзкейс логики
	if cErr := c.uc.Chat(ctx, dConn, chat, uuid); cErr != nil {
		c.log.Error("use case chat", "err", cErr)
	}
}
