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
)

// readPump - Обрабатывает чтение из сокета
func readPump(
	ctx context.Context,
	dataCh chan view.Event,
	cm *decorator.ConnectionManager,
	chat *hub.Hub,
	log *slog.Logger,
) {
	ticker := time.NewTicker(event.PingPeriod)
	defer func() {
		ticker.Stop()
		cm.Close()
	}()

	// ... тут настройка Ping/Pong таймаутов из Кейса 1 ...

	for {
		select {
		case msg := <-dataCh:
			if err := dispatcher(msg, cm, chat, log); err != nil {
				return
			}
		case <-ticker.C:
			// Время слать Ping фронтенду, чтобы проверить, жив ли он
			if err := cm.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				return
			}

			if err := cm.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			// Сигнал от Graceful Shutdown сервера
			if err := cm.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server closed connection"),
			); err != nil {
				return
			}
		}
	}
}

// writePump - Обрабатывает запись в сокет
func writePump(
	ctx context.Context,
	dataCh chan view.Event,
	cm *decorator.ConnectionManager,
	chat *hub.Hub,
	log *slog.Logger,
) {
	ticker := time.NewTicker(event.PingPeriod)
	defer func() {
		ticker.Stop()
		cm.Close()
	}()

	for {
		select {
		case msg := <-dataCh:
			// Пришло сообщение для отправки клиенту
			if err := dispatcher(msg, cm, chat, log); err != nil {
				return
			}
		case <-ticker.C:
			// Время слать Ping фронтенду, чтобы проверить, жив ли он
			if err := cm.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				return
			}

			if err := cm.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			// Сигнал от Graceful Shutdown сервера
			if err := cm.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server closed connection"),
			); err != nil {
				return
			}
		}
	}
}

// dispatcher - Определяет тип события по чтению
func dispatcher(e view.Event, cm *decorator.ConnectionManager, chat *hub.Hub, log *slog.Logger) error {
	switch e.Type {
	case event.JoinEvent:
		var j view.JoinMessagePayload
		if err := json.Unmarshal(e.Payload, &j); err != nil {
			log.Error("failed to unmarshal message payload", "err", err)
		}

		if jErr := chat.Join(cm, j); jErr != nil {
			log.Error("join handler", "err", jErr)
		}
	case event.LeaveEvent:
		var l view.LeaveMessagePayload
		if err := json.Unmarshal(e.Payload, &l); err != nil {
			log.Error("failed to unmarshal message payload", "err", err)
		}

		if lErr := chat.Leave(cm, l); lErr != nil {
			log.Error("leave handler", "err", lErr)
		}
	case event.TextEvent:
		var m view.SendMessagePayload
		if err := json.Unmarshal(e.Payload, &m); err != nil {
			log.Error("failed to unmarshal typing payload", "err", err)
		}

		if tErr := chat.Text(m); tErr != nil {
			log.Error("handler", "err", tErr)
		}
	}

	return nil
}
