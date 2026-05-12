package singleton

import (
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/internal/domain/chat/controller/view"
)

var (
	once     sync.Once
	instance *RoomManager
)

type Room struct {
	Initiator    *Participant  // Инициатор звонка
	Callee       *Participant  // Участник звонка
	PendingOffer *view.SDPData // Данные SDP сессии
}

type Participant struct {
	Conn     *websocket.Conn // Объект коннекшена участника звонка
	Nickname string          // Ник участника звонка
	Role     string          // Роль участника звонка
}

// RoomManager - хранит и управляет комнатами
type RoomManager struct {
	mu    sync.RWMutex
	Rooms map[string]*Room
}

// NewRoom - создает комнату
func NewRoom() *RoomManager {
	once.Do(func() {
		instance = &RoomManager{
			Rooms: make(map[string]*Room),
		}
	})
	return instance
}

// Init - инициализирует комнату с конкретным uuid и инициализирует ее дефолтно
func (rm *RoomManager) Init(uuid string) {
	rm.mu.Lock()
	if _, ok := rm.Rooms[uuid]; !ok {
		rm.Rooms[uuid] = &Room{
			Initiator:    &Participant{},
			Callee:       &Participant{},
			PendingOffer: &view.SDPData{},
		}
	}
	rm.mu.Unlock()
}

// DeleteClient - удаляем коннекшен клиента
func (rm *RoomManager) DeleteClient(uuid string) {
	rm.mu.Lock()
	delete(rm.Rooms, uuid)
	rm.mu.Unlock()
}

// Connections - возвращает количество клиентов
func (rm *RoomManager) Connections() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.Rooms)
}

// Broadcast - рассылает сообщения все клиентам доя установления SDP сессии
func (rm *RoomManager) Broadcast(data view.SDPData, sender *websocket.Conn) error {
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
