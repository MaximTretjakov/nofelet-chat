package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

const (
	joinEvent  = "join"
	textEvent  = "text"
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (uc *UseCase) Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *singleton.Chat) error {
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		if err := cm.Close(); err != nil {
			uc.log.Error("defer close connection chat", "err", err)
		}
	}()

	var e view.Event
	dataCh := make(chan view.Event, 100)

	go writePump(ctx, dataCh, cm, chat, uc.log)
	go readPump(ctx, dataCh, cm, chat, uc.log)

	for {
		if readErr := cm.ReadJSON(&e); readErr != nil {
			uc.log.Error("socket read json", "err", readErr)
			return readErr
		}

		dataCh <- e
	}
}

func readPump(
	ctx context.Context,
	dataCh chan view.Event,
	cm *decorator.ConnectionManager,
	chat *singleton.Chat,
	log *slog.Logger,
) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cm.Close()
	}()

	// ... тут настройка Ping/Pong таймаутов из Кейса 1 ...

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

func writePump(
	ctx context.Context,
	dataCh chan view.Event,
	cm *decorator.ConnectionManager,
	chat *singleton.Chat,
	log *slog.Logger,
) {
	ticker := time.NewTicker(pingPeriod)
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

func dispatcher(e view.Event, cm *decorator.ConnectionManager, chat *singleton.Chat, log *slog.Logger) error {
	switch e.Type {
	case joinEvent:
		var j view.JoinMessagePayload
		if err := json.Unmarshal(e.Payload, &j); err != nil {
			log.Error("failed to unmarshal message payload", "err", err)
		}

		if jErr := join(cm, chat, j); jErr != nil {
			log.Error("join handler", "err", jErr)
		}
	case textEvent:
		var m view.SendMessagePayload
		if err := json.Unmarshal(e.Payload, &m); err != nil {
			log.Error("failed to unmarshal typing payload", "err", err)
		}

		if tErr := text(chat, m); tErr != nil {
			log.Error("handler", "err", tErr)
		}
	}

	return nil
}
