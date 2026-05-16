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
	Conn     *websocket.Conn // Объект коннекшена участника звонка
	Nickname string          // Ник участника звонка
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

// Broadcast - рассылает сообщения все клиентам доя установления SDP сессии
func (rm *Chat) Broadcast(data view.Data, sender *websocket.Conn) error {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for roomID, room := range rm.Rooms {
		initiator := room.Initiator.Conn
		callee := room.Callee.Conn

		if sender == initiator {
			if err := callee.WriteJSON(data); err != nil {
				_ = initiator.Close()
				delete(rm.Rooms, roomID)
				return err
			}
		}

		if err := initiator.WriteJSON(data); err != nil {
			_ = callee.Close()
			delete(rm.Rooms, roomID)
			return err
		}
	}

	return nil
}
