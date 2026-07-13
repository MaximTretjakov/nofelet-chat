package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"nofelet/internal/dependency"
)

func Register(deps *dependency.Container) {
	r := deps.Chat.Routes.Group("/nofelet-chat/api/v1")
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
