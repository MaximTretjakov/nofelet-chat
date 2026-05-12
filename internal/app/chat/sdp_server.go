package chat

import (
	"nofelet/internal/dependency"
	"nofelet/internal/domain/chat"
)

func New(deps *dependency.Container) error {
	chat.Register(deps)

	return nil
}
