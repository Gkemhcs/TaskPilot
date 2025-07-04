package utils

import (
	"github.com/gin-gonic/gin"
	
)

func Success(c *gin.Context, statusCode int,data interface{}) {
	 c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"status":     "error",
		"error_code": statusCode,
		"message":    message,
	})
}
