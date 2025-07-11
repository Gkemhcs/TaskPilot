package utils

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status    string `json:"status" example:"error"`
	ErrorCode int    `json:"error_code" example:"400"`
	Message   string `json:"message" example:"missing email"`
}

// Success sends a JSON response with a success status and the provided data.
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

// Error sends a JSON error response and aborts the request.
func Error(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"status":     "error",
		"error_code": statusCode,
		"message":    message,
	})
}
