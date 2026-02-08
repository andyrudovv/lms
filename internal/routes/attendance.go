package routes

import (
	"net/http"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	DB *db.PostgresConnector
}

// Attendance registers the attendance routes and initializes the handler
func Attendance(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &AttendanceHandler{DB: database}

	attendance := rg.Group("/attendance")
	{
		attendance.POST("", h.markAttendance)
		attendance.GET("/course/:course_id", h.getAttendanceByCourse)
		attendance.GET("/student/:student_id", h.getAttendanceByStudent)
	}
}

// markAttendance handles recording or updating a student's status for a specific date
func (h *AttendanceHandler) markAttendance(c *gin.Context) {
	var input struct {
		CourseID  int    `json:"course_id" binding:"required"`
		StudentID int    `json:"student_id" binding:"required"`
		Date      string `json:"date" binding:"required"` // Format: YYYY-MM-DD
		Status    string `json:"status" binding:"required"` // e.g., 'present', 'absent', 'late'
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SQL Logic: Insert attendance or update status if it already exists for that day (Upsert)
	// Uses the unique index: idx_attendance_course_student_date
	query := `
		INSERT INTO attendance (course_id, student_id, date, status, notes) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_id, student_id, date) 
		DO UPDATE SET status = EXCLUDED.status, notes = EXCLUDED.notes, updated_at = now()`
	
	_, err := h.DB.Exec(c.Request.Context(), query, input.CourseID, input.StudentID, input.Date, input.Status, input.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance recorded successfully"})
}

// getAttendanceByCourse fetches the attendance list for an entire class
func (h *AttendanceHandler) getAttendanceByCourse(c *gin.Context) {
	courseID := c.Param("course_id")

	query := `
		SELECT a.id, u.first_name, u.last_name, a.date, a.status, a.notes 
		FROM attendance a
		JOIN users u ON a.student_id = u.id
		WHERE a.course_id = $1
		ORDER BY a.date DESC`

	rows, err := h.DB.QueryContext(c.Request.Context(), query, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course attendance"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var fname, lname, date, status, notes string
		rows.Scan(&id, &fname, &lname, &date, &status, &notes)
		results = append(results, map[string]interface{}{
			"id":         id,
			"student":    fname + " " + lname,
			"date":       date,
			"status":     status,
			"notes":      notes,
		})
	}

	c.JSON(http.StatusOK, results)
}

// getAttendanceByStudent fetches the attendance history for a specific student
func (h *AttendanceHandler) getAttendanceByStudent(c *gin.Context) {
	studentID := c.Param("student_id")

	query := `
		SELECT a.date, c.title, a.status, a.notes 
		FROM attendance a
		JOIN courses c ON a.course_id = c.id
		WHERE a.student_id = $1
		ORDER BY a.date DESC`

	rows, err := h.DB.QueryContext(c.Request.Context(), query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student attendance history"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var date, title, status, notes string
		rows.Scan(&date, &title, &status, &notes)
		results = append(results, map[string]interface{}{
			"date":   date,
			"course": title,
			"status": status,
			"notes":  notes,
		})
	}

	c.JSON(http.StatusOK, results)
}