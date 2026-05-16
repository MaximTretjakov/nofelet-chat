package usecase

import (
	"errors"

	"nofelet/decorator"
	"nofelet/pkg/singleton"
)

var errChatRoomNotFound = errors.New("chat room not found")

// join - Обрабатывает событие join присоединение к чату
func join(cm *decorator.ConnectionManager, chat *singleton.Chat, uuid string, nick string) error {
	// находим комнату с uuid
	room, ok := chat.Rooms[uuid]
	if !ok {
		return errChatRoomNotFound
	}

	// если комната есть регаемся в ней
	room = append(room, &singleton.ChatRoom{
		Conn:     cm.Conn,
		Nickname: nick,
	})

	return nil
}

// text - Обрабатывает событие text (текстовый тип сообщения)
func text() error {
	return nil
}
