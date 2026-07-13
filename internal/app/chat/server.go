package chat

import (
	"nofelet/internal/dependency"
	"nofelet/internal/domain/chat"
	"nofelet/internal/domain/metrics"
)

func New(deps *dependency.Container) error {
	chat.Register(deps)
	metrics.Register(deps)

	return nil
}
