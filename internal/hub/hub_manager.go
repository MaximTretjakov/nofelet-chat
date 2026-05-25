package hub

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
)

var (
	once     sync.Once
	instance *Hub

	errChatRoomNotFound        = errors.New("chat room not found")
	errChatOrRecipientNotFound = errors.New("chat room or recipient not found")
)

// Hub - управляет чат комнатами
type Hub struct {
	mu    sync.RWMutex
	Rooms map[string]*Room // key = roomID
}

// Room - управляет конкретной комнатой чата
type Room struct {
	mu      sync.RWMutex
	Clients map[string]*Client // key = userID, value = инфо о клиенте
}

// Client - информация о клиенте
type Client struct {
	ID     string
	RoomID string
	Nick   string      // Ник участника звонка
	Send   chan []byte // Канал для writePump
	Conn   *websocket.Conn
}

// NewChatRoom - создает чат комнату
func NewChatRoom() *Hub {
	once.Do(func() {
		instance = &Hub{
			Rooms: make(map[string]*Room),
		}
	})
	return instance
}

// Init - инициализирует комнату с конкретным uuid и инициализирует ее дефолтно
func (h *Hub) Init(uuid string) {
	h.mu.Lock()
	if _, ok := h.Rooms[uuid]; !ok {
		h.Rooms[uuid] = &Room{
			Clients: make(map[string]*Client),
		}
	}
	h.mu.Unlock()
}

// GetRecipient - возвращает пользователя которому адресовано сообщение если он есть в комнате
func (h *Hub) GetRecipient(m view.SendMessagePayload) *Client {
	room, ok := h.Rooms[m.ChatID]
	if !ok {
		return nil
	}

	for _, client := range room.Clients {
		if client.Nick == m.Recipient {
			return client
		}
	}

	return nil
}

// Join - Обрабатывает событие join присоединение к чату
func (h *Hub) Join(cm *decorator.ConnectionManager, j view.JoinMessagePayload) error {
	// находим комнату с uuid
	room, ok := h.Rooms[j.ChatID]
	if !ok {
		return errChatRoomNotFound
	}

	// если комната есть регаемся в ней
	room.mu.Lock()
	room.Clients[j.Nick] = &Client{
		Nick: j.Nick,
		Conn: cm.Conn,
	}
	room.mu.Unlock()

	// todo Шлет user_joined всем остальным в комнате.
	// todo Шлет room_state новому участнику.

	return nil
}

// Leave - Обрабатывает событие leave пользователя
func (h *Hub) Leave(l view.LeaveMessagePayload) error {
	// находим комнату с uuid
	room, ok := h.Rooms[l.ChatID]
	if !ok {
		return errChatRoomNotFound
	}

	room.mu.Lock()
	delete(room.Clients, l.Nick)
	room.mu.Unlock()

	// todo Шлет user_left всем остальным в комнате.

	return nil
}

// Broadcast - бродкастит события всем в чате
func (h *Hub) Broadcast() {}

// Text - Обрабатывает событие text (текстовый тип сообщения)
func (h *Hub) Text(m view.SendMessagePayload) error {
	recipient := h.GetRecipient(m)
	if recipient == nil {
		return errChatOrRecipientNotFound
	}

	if err := recipient.Conn.WriteJSON(m); err != nil {
		return err
	}

	return nil
}
