package routes

import (
	"net/http"
	"time"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/andyrudovv/lms/internal/models"
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

func CreateCourse(db *db.PostgresConnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newCourse models.Course
		if err := c.ShouldBindJSON(&newCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO courses (title, description) VALUES ($1, $2) RETURNING id`
		err := db.QueryRowContext(c.Request.Context(), query, newCourse.Title, newCourse.Description).Scan(&newCourse.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save course"})
			return
		}

		go func(courseTitle string) {
			time.Sleep(2 * time.Second)
			// In a real app, you'd trigger an email or log service here
			println("Background Task: Notifications sent for course:", courseTitle)
		}(newCourse.Title)

		c.JSON(http.StatusCreated, newCourse)
	}
}
