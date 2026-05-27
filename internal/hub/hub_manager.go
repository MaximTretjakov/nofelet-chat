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
		Conn: cm.Conn,
	}
	room.mu.Unlock()

	// шлет JoinEvent всем остальным в комнате
	if bErr := h.Broadcast(
		WithEventType(event.UserJoinedEvent),
		WithJoinMessagePayload(j),
		WithWSConnection(cm.Conn),
	); bErr != nil {
		return bErr
	}

	// создает событие RoomStateEvent
	roomStateEvent, eErr := createRoomStateEvent(room)
	if eErr != nil {
		return eErr
	}

	// шлет RoomStateEvent новому участнику.
	if cErr := cm.WriteJSON(roomStateEvent); cErr != nil {
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
		for _, client := range room.Clients {
			if mt.ws != client.Conn {
				if err := client.Conn.WriteJSON(mt.JoinMessagePayload); err != nil { // todo посылать через Event struct, а не просто payload
					return err
				}
			}
		}
	case event.UserLeftEvent:
		chatID := mt.LeaveMessagePayload.ChatID
		room := h.Rooms[chatID]
		for _, client := range room.Clients {
			if mt.ws != client.Conn {
				if err := client.Conn.WriteJSON(mt.LeaveMessagePayload); err != nil { // todo посылать через Event struct, а не просто payload
					return err
				}
			}
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

	if err := recipient.Conn.WriteJSON(m); err != nil { // todo посылать через Event struct, а не просто SendMessagePayload
		return err
	}

	return nil
}

// createRoomStateEvent - создает событие RoomStateEvent
func createRoomStateEvent(room *Room) (view.Event, error) {
	counter := 1
	clients := make(map[int]string, len(room.Clients))

	for _, client := range room.Clients {
		clients[counter] = client.Nick
		counter++
	}

	payloadData := map[string]interface{}{
		"users": clients,
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return view.Event{}, err
	}

	return view.Event{
		Type:    event.RoomStateEvent,
		Payload: payloadBytes,
	}, nil
}
