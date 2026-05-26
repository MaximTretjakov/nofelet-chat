package view

import (
	"encoding/json"
)

// Event - враппер
type Event struct {
	Type    string          `json:"type"`    // Тип события: "message", "typing", "file"
	Payload json.RawMessage `json:"payload"` // Сырые данные, которые мы распарсим позже
}

// SendMessagePayload - обыкновенное текстовое сообщение
type SendMessagePayload struct {
	ChatID    string `json:"chat_id"`   // ID комнаты чата
	Sender    string `json:"sender"`    // От кого
	Recipient string `json:"recipient"` // Кому
	Content   string `json:"content"`   // Текстовое сообщение
}

// FileMessagePayload - передача файла
type FileMessagePayload struct {
	ChatID   string `json:"chat_id"`   // ID комнаты чата
	FileURL  string `json:"file_url"`  // URL в s3
	FileName string `json:"file_name"` // Имя файла
	FileType string `json:"file_type"` // Типы (image, pdf, etc)
}

// TypingPayload - для анимации набора текста собеседником
type TypingPayload struct {
	ChatID   string `json:"chat_id"`   // ID комнаты чата
	IsTyping bool   `json:"is_typing"` // Флаг того что пользователь что-то набирает в данный момент
}

// JoinMessagePayload - сообщение о присоединении пользователя
type JoinMessagePayload struct {
	ChatID string `json:"chat_id"` // ID комнаты чата
	Nick   string `json:"nick"`    // Ник пользователя
}

// LeaveMessagePayload - сообщение о выходе/дисконнекте пользователя
type LeaveMessagePayload struct {
	JoinMessagePayload
}

// RoomStateMessagePayload - активные пользователи на данный момент
type RoomStateMessagePayload struct {
}
