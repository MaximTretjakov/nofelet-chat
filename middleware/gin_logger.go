package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// DurationLoggerMiddleware - это middleware для логирования запроса
func DurationLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()

		fmt.Printf("Request: Method=%s | Status=%d | Duration=%v | Path=%s\n",
			c.Request.Method,
			status,
			latency,
			c.Request.URL.Path,
		)
	}
}
