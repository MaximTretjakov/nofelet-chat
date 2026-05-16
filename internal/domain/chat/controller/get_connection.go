package controller

import (
	"github.com/gin-gonic/gin"

	"nofelet/decorator"
	"nofelet/pkg/singleton"
)

// GetChat - /chat создает чат
func (c *Controller) GetChat(ctx *gin.Context) {
	uConn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.log.Error("socket creation", "err", sErr)
	}

	chat := singleton.NewChatRoom()

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
