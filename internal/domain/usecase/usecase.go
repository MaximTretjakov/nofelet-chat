package usecase

import (
	"log/slog"

	"nofelet/config"
)

type UseCase struct {
	log *slog.Logger
	cfg *config.Config
}

func New(cfg *config.Config, log *slog.Logger) *UseCase {
	return &UseCase{
		log: log,
		cfg: cfg,
	}
}
