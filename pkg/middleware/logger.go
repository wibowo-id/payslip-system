package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logrus.WithFields(logrus.Fields{
			"status":     status,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"latency":    latency,
			"ip":         c.ClientIP(),
			"request_id": c.GetString("request_id"),
		}).Info("HTTP Request")
	}
}
