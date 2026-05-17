package singleton

import (
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/internal/domain/chat/controller/view"
)

var (
	once     sync.Once
	instance *Chat
)

type ChatRoom struct {
	Conn *websocket.Conn // Объект коннекшена участника звонка
	Nick string          // Ник участника звонка
}

// Chat - управляет чат комнатами
type Chat struct {
	mu    sync.RWMutex
	Rooms map[string][]*ChatRoom
}

// NewChatRoom - создает чат комнату
func NewChatRoom() *Chat {
	once.Do(func() {
		instance = &Chat{
			Rooms: make(map[string][]*ChatRoom),
		}
	})
	return instance
}

// Init - инициализирует комнату с конкретным uuid и инициализирует ее дефолтно
func (rm *Chat) Init(uuid string) {
	rm.mu.Lock()
	if _, ok := rm.Rooms[uuid]; !ok {
		rm.Rooms[uuid] = make([]*ChatRoom, 0)
	}
	rm.mu.Unlock()
}

// DeleteClient - удаляем коннекшен клиента
func (rm *Chat) DeleteClient(uuid string) {
	rm.mu.Lock()
	delete(rm.Rooms, uuid)
	rm.mu.Unlock()
}

// Connections - возвращает количество клиентов
func (rm *Chat) Connections() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.Rooms)
}

// GetUserByName - возвращает пользователя которому адресовано сообщение если он есть в комнате
func (rm *Chat) GetUserByName(m view.SendMessagePayload) *ChatRoom {
	room, ok := rm.Rooms[m.ChatID]
	if !ok {
		return nil
	}

	for _, chatRoom := range room {
		if chatRoom.Nick == m.Recipient {
			return chatRoom
		}
	}

	return nil
}
