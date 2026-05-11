package singleton

import (
	"testing"

	"github.com/google/uuid"
)

func TestRoomManager_LenConnections(t *testing.T) {
	room := NewRoom()
	roomID := uuid.New().String()

	room.Init(roomID)

	tests := []struct {
		name        string
		connections int
		expected    int
	}{
		{
			name:        "проверка количества активных соединений",
			connections: 1,
			expected:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := room.Connections()

			if got != tt.expected {
				t.Errorf("Connections(%d) = %d; want %d", tt.connections, got, tt.expected)
			}
		})
	}
}

func TestRoomManager_DeleteClient(t *testing.T) {
	room := NewRoom()
	roomID := uuid.New().String()

	room.Init(roomID)

	tests := []struct {
		name        string
		roomID      string
		connections int
		expected    int
	}{
		{
			name:        "проверка удаления соединения",
			roomID:      roomID,
			connections: 0,
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room.DeleteClient(tt.roomID)
			got := room.Connections()

			if got != tt.expected {
				t.Errorf("DeleteClient(%d) = %d; want %d", tt.connections, got, tt.expected)
			}
		})
	}
}
