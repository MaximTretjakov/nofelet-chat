package controller

import (
	"log/slog"

	"nofelet/config"
)

type Controller struct {
	log *slog.Logger
	cfg *config.Config
}

func NewController(logger *slog.Logger, cfg *config.Config) *Controller {
	return &Controller{
		log: logger,
		cfg: cfg,
	}
}
