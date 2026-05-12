package chat

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"nofelet/config"
	"nofelet/middleware"
)

type Container struct {
	Routes *gin.Engine
	Logger *slog.Logger
	Cfg    *config.Config
}

// New - создает чат
func New(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	routes, err := newRoutes()
	if err != nil {
		return nil, fmt.Errorf("инициализация роутера: %w", err)
	}

	return &Container{
		Routes: routes,
		Logger: logger,
		Cfg:    cfg,
	}, nil
}

// newRoutes - инициализация роутера
func newRoutes() (*gin.Engine, error) {
	router := gin.New()
	router.ContextWithFallback = true
	router.HandleMethodNotAllowed = true
	router.Use(
		gin.Recovery(),
		middleware.DurationLoggerMiddleware(),
	)
	return router, nil
}
