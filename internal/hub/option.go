package hub

import (
	"github.com/gorilla/websocket"

	"nofelet/internal/domain/chat/controller/view"
)

type MessageTypes struct {
	ws                       *websocket.Conn
	EventType                string
	TypingPayload            view.TypingPayload
	JoinMessagePayload       view.JoinMessagePayload
	LeaveMessagePayload      view.LeaveMessagePayload
	SendMessagePayload       view.SendMessagePayload
	FileMessagePayload       view.FileMessagePayload
	UserJoinedMessagePayload view.UserJoinedMessagePayload
}

func NewMessageTypes() *MessageTypes {
	return &MessageTypes{}
}

type Option func(mt *MessageTypes)

func WithEventType(event string) Option {
	return func(mt *MessageTypes) {
		mt.EventType = event
	}
}

func WithWSConnection(ws *websocket.Conn) Option {
	return func(mt *MessageTypes) {
		mt.ws = ws
	}
}

func WithTypingPayload(payload view.TypingPayload) Option {
	return func(mt *MessageTypes) {
		mt.TypingPayload = payload
	}
}

func WithJoinMessagePayload(payload view.JoinMessagePayload) Option {
	return func(mt *MessageTypes) {
		mt.JoinMessagePayload = payload
	}
}

func WithLeaveMessagePayload(payload view.LeaveMessagePayload) Option {
	return func(mt *MessageTypes) {
		mt.LeaveMessagePayload = payload
	}
}

func WithSendMessagePayload(payload view.SendMessagePayload) Option {
	return func(mt *MessageTypes) {
		mt.SendMessagePayload = payload
	}
}

func WithFileMessagePayload(payload view.FileMessagePayload) Option {
	return func(mt *MessageTypes) {
		mt.FileMessagePayload = payload
	}
}

func WithUserJoinedMessagePayload(payload view.UserJoinedMessagePayload) Option {
	return func(mt *MessageTypes) {
		mt.UserJoinedMessagePayload = payload
	}
}
