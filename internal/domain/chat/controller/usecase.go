package controller

import (
	"context"

	"nofelet/decorator"
	"nofelet/internal/hub"
)

type UseCase interface {
	// Chat - логика чата
	Chat(ctx context.Context, cm *decorator.ConnectionManager, chat *hub.Hub) error
}
