package routes

import (
	"net/http"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	DB *db.PostgresConnector
}

// Progress registers the progress-related routes
func Progress(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &ProgressHandler{DB: database}

	// Route: GET /api/v1/courses/:id/progress?student_id=123
	rg.GET("/courses/:id/progress", h.GetProgress)
	
	// Route: GET /api/v1/students/:student_id/overall-progress
	rg.GET("/students/:student_id/overall-progress", h.GetOverallProgress)
}

// GetProgress calculates the completion percentage for a specific course
func (h *ProgressHandler) GetProgress(c *gin.Context) {
	studentID := c.Query("student_id")
	courseID := c.Param("id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id query parameter is required"})
		return
	}

	// SQL Logic:
	// 1. Total weeks defined for the course
	// 2. Count of unique dates in attendance where student was 'present'
	query := `
		SELECT 
			(SELECT COUNT(id) FROM weeks WHERE course_id = $1) as total_weeks,
			(SELECT COUNT(DISTINCT date) FROM attendance 
			 WHERE course_id = $1 AND student_id = $2 AND status = 'present') as attended_days
	`

	var total, attended int
	err := h.DB.QueryRowContext(c.Request.Context(), query, courseID, studentID).Scan(&total, &attended)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate progress"})
		return
	}

	// Handle division by zero if course has no weeks yet
	var percent float64
	if total > 0 {
		percent = (float64(attended) / float64(total)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"course_id":             courseID,
		"student_id":            studentID,
		"total_modules":        total,
		"completed_modules":    attended,
		"completion_percentage": percent,
	})
}

// GetOverallProgress fetches progress for all courses a student is enrolled in
func (h *ProgressHandler) GetOverallProgress(c *gin.Context) {
	studentID := c.Param("student_id")

	query := `
		SELECT 
			c.id, 
			c.title,
			(SELECT COUNT(w.id) FROM weeks w WHERE w.course_id = c.id) as total,
			(SELECT COUNT(DISTINCT a.date) FROM attendance a 
			 WHERE a.course_id = c.id AND a.student_id = e.user_id AND a.status = 'present') as attended
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE e.user_id = $1
	`

	rows, err := h.DB.QueryContext(c.Request.Context(), query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch overall progress"})
		return
	}
	defer rows.Close()

	type CourseProgress struct {
		CourseID   int     `json:"course_id"`
		Title      string  `json:"title"`
		Percentage float64 `json:"percentage"`
	}

	var results []CourseProgress
	for rows.Next() {
		var id int
		var title string
		var total, attended int
		rows.Scan(&id, &title, &total, &attended)

		percent := 0.0
		if total > 0 {
			percent = (float64(attended) / float64(total)) * 100
		}

		results = append(results, CourseProgress{
			CourseID:   id,
			Title:      title,
			Percentage: percent,
		})
	}

	c.JSON(http.StatusOK, results)
}