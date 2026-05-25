package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"nofelet/decorator"
	"nofelet/internal/hub"
)

var errChatRoomNotFound = errors.New("chat room not found")

// GetChat - /chat создает чат
func (c *Controller) GetChat(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if len(uuid) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errChatRoomNotFound)
		return
	}

	uConn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.log.Error("socket creation", "err", sErr)
	}

	chat := hub.NewChatRoom()
	chat.Init(uuid)

	dConn := decorator.New(uConn, c.log, c.cfg)
	defer func() {
		if err := dConn.Close(); err != nil {
			c.log.Error("close connection", "err", err)
		}
	}()

	if cErr := c.uc.Chat(ctx, dConn, chat); cErr != nil {
		c.log.Error("use case chat", "err", cErr)
	}
}
