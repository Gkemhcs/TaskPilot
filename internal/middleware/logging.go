package middleware

import (
	
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware (logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming request")

		
		
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		err := c.Errors.Last()

		entry := logger.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": latency,
		})
		if err != nil || status >= 400 {
    if err != nil {
        entry.WithField("error", err.Err).Error("Request failed")
    } else {
        entry.Error("Request failed with status code")
    }
} else {
			entry.Info("Request completed")
		}

	}

}
