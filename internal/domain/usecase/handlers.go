package usecase

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errChatRoomNotFound        = errors.New("chat room not found")
	errChatOrRecipientNotFound = errors.New("chat room or recipient not found")
)

// join - Обрабатывает событие join присоединение к чату
func join(cm *decorator.ConnectionManager, chat *singleton.Chat, j view.JoinMessagePayload) error {
	// находим комнату с uuid
	room, ok := chat.Rooms[j.ChatID]
	if !ok {
		return errChatRoomNotFound
	}

	// если комната есть регаемся в ней
	room = append(room, &singleton.ChatRoom{
		Conn: cm.Conn,
		Nick: j.Nick,
	})

	return nil
}

// text - Обрабатывает событие text (текстовый тип сообщения)
func text(chat *singleton.Chat, m view.SendMessagePayload) error {
	recipient := chat.GetUserByName(m)
	if recipient == nil {
		return errChatOrRecipientNotFound
	}

	if err := recipient.Conn.WriteJSON(m); err != nil {
		return err
	}

	return nil
}

// readPump - Обрабатывает чтение из сокета
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
			if err := readDispatcher(msg, cm, chat, log); err != nil {
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
			if err := writeDispatcher(msg, chat, log); err != nil {
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

// readDispatcher - Определяет тип события по чтению
func readDispatcher(e view.Event, cm *decorator.ConnectionManager, chat *singleton.Chat, log *slog.Logger) error {
	switch e.Type {
	case joinEvent:
		var j view.JoinMessagePayload
		if err := json.Unmarshal(e.Payload, &j); err != nil {
			log.Error("failed to unmarshal message payload", "err", err)
		}

		if jErr := join(cm, chat, j); jErr != nil {
			log.Error("join handler", "err", jErr)
		}
	}

	return nil
}

// writeDispatcher - Определяет тип события по записи
func writeDispatcher(e view.Event, chat *singleton.Chat, log *slog.Logger) error {
	switch e.Type {
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
