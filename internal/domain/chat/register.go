package chat

import (
	"nofelet/internal/dependency"
	sc "nofelet/internal/domain/chat/controller"
)

func Register(deps *dependency.Container) {
	c := sc.NewController(deps.Logger, deps.Cfg)

	r := deps.Chat.Routes
	r.GET("/connect/:uuid", c.GetConnection)
}
