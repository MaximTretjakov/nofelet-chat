package controller

import (
	"log/slog"

	"nofelet/config"
)

type Controller struct {
	Logger *slog.Logger
	Config *config.Config
}

func NewController(logger *slog.Logger, cfg *config.Config) *Controller {
	return &Controller{
		Logger: logger,
		Config: cfg,
	}
}
