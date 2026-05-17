package usecase

import (
	"errors"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
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
