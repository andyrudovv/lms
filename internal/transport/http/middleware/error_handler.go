package middleware

import (
	"net/http"

	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors[0].Err
		responder.Fail(c, http.StatusBadRequest, err.Error())
	}
}
