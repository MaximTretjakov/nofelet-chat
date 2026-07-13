package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Активные запросы
		activeRequests.Add(c, 1)
		defer activeRequests.Add(c, -1)

		c.Next()

		// Счетчик всех HTTP-запросов
		requestsTotal.Add(c.Request.Context(), 1)

		// Гистограмма времени выполнения запроса
		requestDuration.Record(c.Request.Context(), time.Since(start).Seconds())

		// Счетчик status code
		responseTotal.Add(
			c,
			1,
			metric.WithAttributes(
				attribute.Int("status", c.Writer.Status()),
			),
		)
	}
}
