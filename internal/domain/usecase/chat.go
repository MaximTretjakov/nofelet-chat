package usecase

import (
	"context"
	"errors"
	"net"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/internal/hub"
)

var errClientDisconnected = errors.New("client is disconnected")

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
		if rErr := cm.ReadJSON(&e); rErr != nil {
			// Проверяем, является ли ошибка таймаутом
			var netErr net.Error
			if errors.As(rErr, &netErr) && netErr.Timeout() {
				uc.log.Error("use case chat", "err", errClientDisconnected)
			}
			uc.log.Error("use case chat", "err", rErr)
			return rErr
		}

		dataCh <- e
	}
}
