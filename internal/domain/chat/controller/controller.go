package controller

import (
	"log/slog"

	"nofelet/config"
)

type Controller struct {
	uc  UseCase
	log *slog.Logger
	cfg *config.Config
}

func New(uc UseCase, log *slog.Logger, cfg *config.Config) *Controller {
	return &Controller{
		uc:  uc,
		log: log,
		cfg: cfg,
	}
}
