package chat

import (
	"nofelet/internal/dependency"
	sc "nofelet/internal/domain/chat/controller"
	"nofelet/internal/domain/usecase"
)

func Register(deps *dependency.Container) {
	c := sc.New(makeUC(deps), deps.Logger, deps.Cfg)

	r := deps.Chat.Routes
	r.GET("/chat/:uuid", c.GetChat)
}

func makeUC(deps *dependency.Container) *usecase.UseCase {
	return usecase.New(deps.Cfg, deps.Logger)
}
