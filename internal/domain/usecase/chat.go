package usecase

import (
	"context"
	"log/slog"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/internal/hub"
)

func (uc *UseCase) Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *hub.Hub, uuid string) error {
	var e view.Event
	dataCh := make(chan view.Event, 100)
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		close(dataCh)
		if err := cm.Close(); err != nil {
			uc.log.Error("use case chat", "err", err)
		}
	}()

	// Обработка входящих сообщений
	go handler(ctx, dataCh, cm, chat, uc.log)

	// Прием входящих сообщений
	for {
		if rErr := cm.ReadJSON(&e); rErr != nil {
			handleError(rErr, cm, chat, uuid, uc.log)
			return rErr
		}

		dataCh <- e
	}
}

func handleError(err error, cm *decorator.ConnectionManager, chat *hub.Hub, uuid string, log *slog.Logger) {
	// шлем всем остальным клиентам что такой-то клиент отвалился
	l := view.LeaveMessagePayload{
		ChatID: uuid,
	}
	if lErr := chat.Leave(cm, l); lErr != nil {
		log.Error("use case chat", "err", lErr)
	}

	log.Error("use case chat", "err", err)
}
