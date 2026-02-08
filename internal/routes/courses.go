package routes

import (
	"net/http"
	"time"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/andyrudovv/lms/internal/models"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	DB *db.PostgresConnector
}

func Courses(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &CourseHandler{DB: database}
	courses := rg.Group("/courses")
	{
		courses.GET("", h.getCourses)
		courses.POST("", h.createCourse)
		courses.GET("/:id", h.getCourseByID)
	}
}

func (h *CourseHandler) getCourses(c *gin.Context) {
	query := `SELECT id, title, description, teacher_id, created_at FROM courses`
	rows, err := h.DB.QueryContext(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var crs models.Course
		rows.Scan(&crs.ID, &crs.Title, &crs.Description, &crs.TeacherID, &crs.CreatedAt)
		courses = append(courses, crs)
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *CourseHandler) createCourse(c *gin.Context) {
	var newCourse models.Course
	if err := c.ShouldBindJSON(&newCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO courses (title, description, teacher_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := h.DB.QueryRowContext(c.Request.Context(), query, 
		newCourse.Title, newCourse.Description, newCourse.TeacherID).Scan(&newCourse.ID, &newCourse.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save course"})
		return
	}

	// Async Notification/Publishing Log
	go func(title string) {
		time.Sleep(2 * time.Second)
		println("Background Task: Notifications sent for course publishing:", title)
	}(newCourse.Title)

	c.JSON(http.StatusCreated, newCourse)
}

func (h *CourseHandler) getCourseByID(c *gin.Context) {
	id := c.Param("id")
	var crs models.Course
	query := `SELECT id, title, description, teacher_id, created_at FROM courses WHERE id = $1`
	err := h.DB.QueryRowContext(c.Request.Context(), query, id).Scan(&crs.ID, &crs.Title, &crs.Description, &crs.TeacherID, &crs.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, crs)
}