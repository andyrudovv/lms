package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
	}
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "login endpoint",
	})
}

func register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "register endpoint",
	})
}
