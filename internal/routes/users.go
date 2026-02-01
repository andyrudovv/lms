package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Users(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", getUsers)
		users.GET("/:id", getUserByID)
	}
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []string{},
	})
}

func getUserByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": c.Param("id"),
	})
}
