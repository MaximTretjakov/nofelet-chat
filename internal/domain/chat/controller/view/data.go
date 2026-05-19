package view

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`    // Тип события: "message", "typing", "file"
	Payload json.RawMessage `json:"payload"` // Сырые данные, которые мы распарсим позже
}

type SendMessagePayload struct {
	ChatID    string `json:"chat_id"`   // ID комнаты чата
	Sender    string `json:"sender"`    // От кого
	Recipient string `json:"recipient"` // Кому
	Content   string `json:"content"`   // Текстовое сообщение
}

type FileMessagePayload struct {
	ChatID   string `json:"chat_id"`   // ID комнаты чата
	FileURL  string `json:"file_url"`  // URL в s3
	FileName string `json:"file_name"` // Имя файла
	FileType string `json:"file_type"` // Типы (image, pdf, etc)
}

type TypingPayload struct {
	ChatID   string `json:"chat_id"`   // ID комнаты чата
	IsTyping bool   `json:"is_typing"` // Флаг того что пользователь что-то набирает в данный момент
}

type JoinMessagePayload struct {
	ChatID string `json:"chat_id"` // ID комнаты чата
	Nick   string `json:"nick"`    // Ник пользователя
}
