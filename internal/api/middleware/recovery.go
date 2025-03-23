package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the error
				log := logger.WithRequestContext(c.Request.Context())
				log.WithFields(logrus.Fields{
					"error": err,
					"stack": stack,
				}).Error("Panic recovered")

				// Respond with 500 Internal Server Error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
