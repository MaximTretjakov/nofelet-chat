package event

const (
	JoinRoomEvent    = "join_room"    // c->s Отправляется сразу после установки WS-соединения
	SendMessageEvent = "send_message" // c->s Отправляется, когда юзер нажал «Отправить» в чате

	RoomStateEvent  = "room_state"  // s->c Приходит только один раз в ответ на join_room
	UserJoinedEvent = "user_joined" // s->c Приходит всем, когда кто-то новый подключился
	UserLeftEvent   = "leave_room"  // s->c Приходит всем, когда кто-то отключился
	NewMessageEvent = "new_message" // s->c Приходит всем, когда кто-то прислал текстовое сообщение
)
