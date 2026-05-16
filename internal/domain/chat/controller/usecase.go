package controller

import (
	"context"

	"nofelet/decorator"
	"nofelet/pkg/singleton"
)

type UseCase interface {
	// Chat - логика чата
	Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *singleton.Chat, uuid string) error
}
