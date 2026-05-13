package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

// GetChat - /chat создает чат
func (c *Controller) GetChat(ctx *gin.Context) {
	uConn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.log.Error("socket creation", slog.Any("err", sErr))
	}

	chat := singleton.NewChatRoom()

	dConn := decorator.NewConnection(uConn, c.log, c.cfg)
	defer func() {
		if err := dConn.Close(); err != nil {
			c.log.Error("close connection", slog.Any("err", err))
		}
	}()

	go handler(dConn, chat, c.log)
}

// handler - обрабатывает коннекты участников
func handler(mc *decorator.ConnectionManager, chat *singleton.Chat, logger *slog.Logger) {
	defer func() {
		chat.DeleteClient(uuid)
		_ = mc.Close()
	}()

	var data view.Data

	for {
		if readErr := mc.ReadJSON(&data); readErr != nil {
			logger.Error("socket read", slog.Any("err", readErr))
			break
		}

		switch data.Type {
		case "join":
			jErr := Join(data, mc.Conn, r, room)
			if jErr != nil {
				logger.Error("handler", slog.Any("join error:", jErr))
			}
		case "offer":
			oErr := Offer(data, mc.Conn, r, room)
			if oErr != nil {
				logger.Error("handler", slog.Any("offer error:", oErr))
			}
		case "ice-candidate":
			iceErr := IceCandidate(data, mc.Conn, r, room)
			if iceErr != nil {
				logger.Error("handler", slog.Any("ice-candidate error:", iceErr))
			}
		case "answer":
			brErr := room.Broadcast(data, mc.Conn)
			if brErr != nil {
				logger.Error("handler", slog.Any("answer error:", brErr))
			}
		}
	}
}
