package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/internal/domain/event"
	"nofelet/internal/hub"
	"nofelet/middleware/metrics"
)

const (
	delta      = 10 * time.Second // Смещаем время следубщего пинга
	pingPeriod = 10 * time.Second // Как часто слать Ping
)

// handler - Обрабатывает клиентские сообщения
func handler(
	ctx context.Context,
	dataCh chan view.Event,
	cm *decorator.ConnectionManager,
	chat *hub.Hub,
	log *slog.Logger,
) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := cm.Close(); err != nil {
			log.Error("chat handler", "err", err)
		}
	}()

	for {
		select {
		case msg := <-dataCh:
			if err := dispatcher(ctx, msg, cm, chat, log); err != nil {
				return
			}
		case <-ticker.C:
			// Ping фронтенду, жив ли он?
			if err := cm.SetWriteDeadline(time.Now().Add(delta)); err != nil {
				log.Error("chat handler", "err", err)
				return
			}
			if err := cm.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("chat handler", "err", err)
				return
			}
		case <-ctx.Done():
			if err := cm.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server closed connection"),
			); err != nil {
				log.Error("chat handler", "err", err)
				return
			}
		}
	}
}

// dispatcher - Обрабатывает события от клиента
func dispatcher(
	ctx context.Context,
	e view.Event,
	cm *decorator.ConnectionManager,
	chat *hub.Hub,
	log *slog.Logger,
) error {
	switch e.Type {
	case event.JoinRoomEvent:
		var j view.JoinMessagePayload
		if err := json.Unmarshal(e.Payload, &j); err != nil {
			log.Error("failed to unmarshal message payload", "err", err)
		}

		if jErr := chat.Join(cm, j); jErr != nil {
			log.Error("join handler", "err", jErr)
		}
	case event.SendMessageEvent:
		var m view.SendMessagePayload
		if err := json.Unmarshal(e.Payload, &m); err != nil {
			log.Error("failed to unmarshal typing payload", "err", err)
		}

		if tErr := chat.SendMessage(m); tErr != nil {
			log.Error("handler", "err", tErr)
			metrics.MessagesTotal.Add(ctx, -1)
		}
		metrics.MessagesTotal.Add(ctx, 1)
	}

	return nil
}
