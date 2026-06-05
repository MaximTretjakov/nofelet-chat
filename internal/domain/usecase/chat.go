package usecase

import (
	"context"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/internal/hub"
)

func (uc *UseCase) Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *hub.Hub) error {
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
		if readErr := cm.ReadJSON(&e); readErr != nil {
			uc.log.Error("socket read json", "err", readErr)
			return readErr
		}

		dataCh <- e
	}
}
