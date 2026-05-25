package singleton

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"nofelet/internal/domain/chat/controller/view"
)

var (
	once     sync.Once
	instance *Hub
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

// AddClient - добавляет клиента в комнату
func (r *Room) AddClient(nick string, conn *websocket.Conn) string {
	clientID := uuid.New().String()

	r.mu.Lock()
	r.Clients[clientID] = &Client{
		Nick: nick,
		Conn: conn,
	}
	r.mu.Unlock()

	return clientID
}

// RemoveClient - удаляет клиента из комнаты
func (r *Room) RemoveClient(clientID string) {
	r.mu.Lock()
	delete(r.Clients, clientID)
	r.mu.Unlock()
}
