package hub

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/internal/domain/event"
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
	RoomID string
	Nick   string // Ник участника звонка
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
		Conn: cm.Conn, // todo как-то задекорировать чтобы видеть в логах
	}
	room.mu.Unlock()

	// шлет JoinEvent всем остальным в комнате
	if bErr := h.Broadcast(
		WithEventType(event.UserJoinedEvent),
		WithJoinMessagePayload(j),
		WithUserJoinedMessagePayload(view.UserJoinedMessagePayload{Nick: j.Nick}),
		WithWSConnection(cm.Conn),
	); bErr != nil {
		return bErr
	}

	// создает payload для события RoomStateEvent
	payload, pErr := createPayload(room)
	if pErr != nil {
		return pErr
	}

	// создает событие RoomStateEvent
	rsEvent, eErr := newEvent(event.RoomStateEvent, payload)
	if eErr != nil {
		return eErr
	}

	// шлет RoomStateEvent новому участнику.
	if cErr := cm.WriteJSON(rsEvent); cErr != nil {
		return cErr
	}

	return nil
}

// Leave - Обрабатывает событие leave пользователя
func (h *Hub) Leave(cm *decorator.ConnectionManager, l view.LeaveMessagePayload) error {
	// находим комнату с uuid
	room, ok := h.Rooms[l.ChatID]
	if !ok {
		return errChatRoomNotFound
	}

	room.mu.Lock()
	delete(room.Clients, l.Nick)
	room.mu.Unlock()

	// шлет UserLeftEvent всем остальным в комнате
	if bErr := h.Broadcast(
		WithEventType(event.UserLeftEvent),
		WithLeaveMessagePayload(l),
		WithWSConnection(cm.Conn),
	); bErr != nil {
		return bErr
	}

	return nil
}

// Broadcast - бродкастит события всем в чате
func (h *Hub) Broadcast(options ...Option) error {
	mt := NewMessageTypes()

	for _, option := range options {
		option(mt)
	}

	switch mt.EventType {
	case event.UserJoinedEvent:
		chatID := mt.JoinMessagePayload.ChatID
		room := h.Rooms[chatID]

		ujEvent, eErr := newEvent(
			event.UserJoinedEvent,
			map[string]interface{}{"user_joined": mt.UserJoinedMessagePayload},
		)
		if eErr != nil {
			return eErr
		}

		for _, client := range room.Clients {
			if mt.ws != client.Conn {
				room.mu.Lock()
				if err := client.Conn.WriteJSON(ujEvent); err != nil {
					return err
				}
				room.mu.Unlock()
			}
		}

	case event.UserLeftEvent:
		chatID := mt.LeaveMessagePayload.ChatID
		room := h.Rooms[chatID]

		ulEvent, eErr := newEvent(
			event.UserLeftEvent,
			map[string]interface{}{"leave_room": mt.LeaveMessagePayload},
		)
		if eErr != nil {
			return eErr
		}

		for _, client := range room.Clients {
			room.mu.Lock()
			if mt.ws != client.Conn {
				if err := client.Conn.WriteJSON(ulEvent); err != nil {
					return err
				}
			}
			room.mu.Lock()
		}
	}

	return nil
}

// SendMessage - Обрабатывает событие text (текстовый тип сообщения)
func (h *Hub) SendMessage(m view.SendMessagePayload) error {
	recipient := h.GetRecipient(m)
	if recipient == nil {
		return errChatOrRecipientNotFound
	}

	smEvent, eErr := newEvent(
		event.NewMessageEvent,
		map[string]interface{}{"new_message": m},
	)
	if eErr != nil {
		return eErr
	}

	if err := recipient.Conn.WriteJSON(smEvent); err != nil {
		return err
	}

	return nil
}

// createPayload - подготавливает payload для события RoomStateEvent
func createPayload(room *Room) (map[string]interface{}, error) {
	clients := make(map[string]string, len(room.Clients))
	for _, client := range room.Clients {
		clients[client.Nick] = client.Nick
	}

	return map[string]interface{}{"room_state": clients}, nil
}

// newEvent - создает событие
func newEvent(eventType string, data interface{}) (view.Event, error) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return view.Event{}, err
	}

	return view.Event{
		Type:    eventType,
		Payload: payloadBytes,
	}, nil
}
