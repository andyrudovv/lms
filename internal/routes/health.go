package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})
}
