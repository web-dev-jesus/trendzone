package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Create a request context with request information
		requestCtx := logger.NewRequestContext(
			c.Request.Context(),
			c.Request.URL.String(),
			c.ClientIP(),
		)
		c.Request = c.Request.WithContext(requestCtx)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log the request details
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()

		// Get the error if there was one
		err := c.Errors.Last()
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}

		log := logger.WithRequestContext(requestCtx)

		// Customize log based on status code
		if statusCode >= 500 {
			log.WithFields(logrus.Fields{
				"status":     statusCode,
				"latency_ms": latency.Milliseconds(),
				"method":     method,
				"path":       path,
				"query":      query,
				"user_agent": userAgent,
				"error":      errMsg,
			}).Error("Server error")
		} else if statusCode >= 400 {
			log.WithFields(logrus.Fields{
				"status":     statusCode,
				"latency_ms": latency.Milliseconds(),
				"method":     method,
				"path":       path,
				"query":      query,
				"user_agent": userAgent,
				"error":      errMsg,
			}).Warn("Client error")
		} else {
			log.WithFields(logrus.Fields{
				"status":     statusCode,
				"latency_ms": latency.Milliseconds(),
				"method":     method,
				"path":       path,
				"query":      query,
				"user_agent": userAgent,
			}).Info("Request processed")
		}
	}
}
