package controller

import (
	"nofelet/decorator"
	"nofelet/pkg/singleton"
)

type UseCase interface {
	// Chat - логика чата
	Chat(cm *decorator.ConnectionManager, chat *singleton.Chat) error
}
