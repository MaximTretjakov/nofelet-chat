package usecase

import (
	"context"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

func (uc *UseCase) Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *singleton.Chat) error {
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		if err := cm.Close(); err != nil {
			uc.log.Error("defer close connection chat", "err", err)
		}
	}()

	var e view.Event
	dataCh := make(chan view.Event, 100)

	go writePump(ctx, dataCh, cm, chat, uc.log)
	go readPump(ctx, dataCh, cm, chat, uc.log)

	for {
		if readErr := cm.ReadJSON(&e); readErr != nil {
			uc.log.Error("socket read json", "err", readErr)
			return readErr
		}

		dataCh <- e
	}
}
