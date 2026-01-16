package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Courses(rg *gin.RouterGroup) {
	courses := rg.Group("/courses")
	{
		courses.GET("", getCourses)
		courses.POST("", createCourse)
		courses.GET("/:id", getCourseByID)
	}
}

func getCourses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"courses": []string{},
	})
}

func createCourse(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "course created",
	})
}

func getCourseByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": c.Param("id"),
	})
}
