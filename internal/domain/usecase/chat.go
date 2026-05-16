package usecase

import (
	"context"
	"encoding/json"
	"log/slog"

	"nofelet/decorator"
	"nofelet/internal/domain/chat/controller/view"
	"nofelet/pkg/singleton"
)

const (
	joinEvent = "join"
	textEvent = "text"
)

func (uc *UseCase) Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *singleton.Chat, uuid string) error {
	go handler(cm, chat, uc.log, uuid)

	return nil
}

// handler - обрабатывает логику обмена сообщениями
func handler(cm *decorator.ConnectionManager, chat *singleton.Chat, log *slog.Logger, uuid string) {
	var e view.Event

	for {
		if readErr := cm.ReadJSON(&e); readErr != nil {
			log.Error("socket read json", "err", readErr)
			break
		}

		switch e.Type {
		case joinEvent:
			var msg view.JoinMessagePayload
			if err := json.Unmarshal(e.Payload, &msg); err != nil {
				log.Error("failed to unmarshal message payload", "err", err)
				continue // эту ошибку обработали, идем читать следующее сообщение
			}

			if jErr := join(cm, chat, uuid, msg.Nick); jErr != nil {
				log.Error("join handler", "err", jErr)
			}
		case textEvent:
			var typing view.TypingPayload
			if err := json.Unmarshal(e.Payload, &typing); err != nil {
				log.Error("failed to unmarshal typing payload", "err", err)
				continue
			}

			if tErr := text(); tErr != nil {
				log.Error("handler", "err", tErr)
			}
		}
	}
}
