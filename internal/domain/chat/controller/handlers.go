package controller

import (
	"github.com/gorilla/websocket"

	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

const (
	caller = "caller"
	callee = "callee"
)

// Join - Обрабатывает событие join
func Join(data view.SDPData, conn *websocket.Conn, r *singleton.Room, rs *singleton.RoomManager) error {
	// Кто пришел? Инициатор или собеседник?
	if data.Participant.Role == caller {
		r.Initiator.Conn = conn
		r.Initiator.Nickname = data.Participant.Name
		r.Initiator.Role = data.Participant.Role
	}
	if data.Participant.Role == callee {
		r.Callee.Conn = conn
		r.Callee.Nickname = data.Participant.Name
		r.Callee.Role = data.Participant.Role
	}
	// Если инициатор и собеседник на месте то начинаем обмен никами и возможно офером если он есть
	if r.Initiator.Conn != nil && r.Callee.Conn != nil {
		// Шлем инициатору ник собеседника
		r.Initiator.Conn.WriteJSON(view.SDPData{
			Type: "ready",
			Participant: view.Participant{
				Name: r.Callee.Nickname,
			},
		})
		// Шлем собеседнику ник инициатора
		r.Callee.Conn.WriteJSON(view.SDPData{
			Type: "ready",
			Participant: view.Participant{
				Name: r.Initiator.Nickname,
			},
		})
		// ПУНКТ 7: Если был припасен оффер — сразу отдаем его Callee
		if r.PendingOffer != nil {
			if brErr := rs.Broadcast(*r.PendingOffer, r.Initiator.Conn); brErr != nil {
				return brErr
			}
			r.PendingOffer = nil
		}
	}

	return nil
}

// Offer - Обрабатывает событие offer
func Offer(data view.SDPData, conn *websocket.Conn, r *singleton.Room, rs *singleton.RoomManager) error {
	if r.Callee.Conn == nil {
		offer := data
		r.PendingOffer = &offer
		return nil
	}

	if brErr := rs.Broadcast(data, conn); brErr != nil {
		return brErr
	}

	return nil
}

// IceCandidate - Обрабатывает событие ice-candidate
func IceCandidate(data view.SDPData, conn *websocket.Conn, r *singleton.Room, rs *singleton.RoomManager) error {
	data.SDP = ""

	if r.Callee.Conn == nil {
		return nil
	}

	if brErr := rs.Broadcast(data, conn); brErr != nil {
		return brErr
	}

	return nil
}
