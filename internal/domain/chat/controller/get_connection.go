package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

// GetConnection - /connect/:uuid установка sdp сессии
func (c *Controller) GetConnection(ctx *gin.Context) {
	upgradedConn, sErr := Upgrader(ctx)
	if sErr != nil {
		c.Logger.Error("socket creation", slog.Any("err", sErr))
	}

	managedConn := decorator.NewConn(upgradedConn, c.Logger, c.Config)

	go handler(managedConn, uuid, room, c.Logger)
}

// handler - обрабатывает коннекты участников
func handler(mc *decorator.ConnectionManager, uuid string, room *singleton.RoomManager, logger *slog.Logger) {
	defer func() {
		room.DeleteClient(uuid)
		_ = mc.Close()
	}()

	var data view.SDPData

	for {
		if readErr := mc.ReadJSON(&data); readErr != nil {
			logger.Error("socket read", slog.Any("err", readErr))
			break
		}

		r := room.Rooms[uuid]

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
