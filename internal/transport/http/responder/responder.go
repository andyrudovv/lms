package responder

import "github.com/gin-gonic/gin"

type APIError struct {
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, gin.H{"data": data})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, gin.H{"data": data})
}

func Fail(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"error": APIError{Message: msg}})
}
