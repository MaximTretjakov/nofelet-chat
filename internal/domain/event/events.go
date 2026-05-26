package event

import "time"

const (
	JoinEvent      = "join"
	LeaveEvent     = "leave"
	TextEvent      = "text"
	RoomStateEvent = "roomState"
	PongWait       = 60 * time.Second
	PingPeriod     = (PongWait * 9) / 10
)
