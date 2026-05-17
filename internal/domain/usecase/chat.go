package usecase

import (
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

func (uc *UseCase) Chat(cm *decorator.ConnectionManager, chat *singleton.Chat) error {
	go handler(cm, chat, uc.log) // todo реализовать контроль за горутинами (в каких случаях оа завершается и как)

	return nil
}

// handler - обрабатывает логику обмена сообщениями
func handler(cm *decorator.ConnectionManager, chat *singleton.Chat, log *slog.Logger) {
	var e view.Event

	for {
		if readErr := cm.ReadJSON(&e); readErr != nil {
			log.Error("socket read json", "err", readErr)
			break
		}

		switch e.Type {
		case joinEvent:
			var j view.JoinMessagePayload
			if err := json.Unmarshal(e.Payload, &j); err != nil {
				log.Error("failed to unmarshal message payload", "err", err)
				continue // эту ошибку обработали, идем читать следующее сообщение
			}

			if jErr := join(cm, chat, j); jErr != nil {
				log.Error("join handler", "err", jErr)
			}
		case textEvent:
			var m view.SendMessagePayload
			if err := json.Unmarshal(e.Payload, &m); err != nil {
				log.Error("failed to unmarshal typing payload", "err", err)
				continue
			}

			if tErr := text(chat, m); tErr != nil {
				log.Error("handler", "err", tErr)
			}
		}
	}
}
