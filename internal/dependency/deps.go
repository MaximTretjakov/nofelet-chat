package dependency

import (
	"fmt"
	"log/slog"

	"nofelet/config"
	"nofelet/internal/dependency/chat"
)

// Container - основной контейнер зависимостей
type Container struct {
	Chat   *chat.Container
	Logger *slog.Logger
	Cfg    *config.Config
}

// New - создает DI контейнер
func New(Cfg *config.Config, logger *slog.Logger) (*Container, error) {
	ChatContainer, err := chat.New(Cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("создание контейнера di: %w", err)
	}

	return &Container{
		Chat:   ChatContainer,
		Logger: logger,
		Cfg:    Cfg,
	}, nil
}
