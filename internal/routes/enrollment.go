package routes

import (
	"net/http"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	DB *db.PostgresConnector
}

// Enrollment registers the enrollment routes
func Enrollment(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &EnrollmentHandler{DB: database}

	enrollments := rg.Group("/enrollments")
	{
		enrollments.POST("", h.EnrollStudent)             // Logic for enrolling a student
		enrollments.GET("/user/:user_id", h.GetByUser)    // View student's courses
		enrollments.GET("/course/:course_id", h.GetByCourse) // View course's roster
		enrollments.DELETE("/:id", h.Unenroll)           // Remove enrollment
	}
}

// EnrollStudent handles the logic for student enrollment
func (h *EnrollmentHandler) EnrollStudent(c *gin.Context) {
	var req struct {
		UserID   int `json:"user_id" binding:"required"`
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SQL insertion using the unique constraint idx_enrollments_user_course from your schema
	query := `INSERT INTO enrollments (user_id, course_id) VALUES ($1, $2)`
	
	_, err := h.DB.Exec(c.Request.Context(), query, req.UserID, req.CourseID)
	if err != nil {
		// Typically a 409 Conflict if the student is already enrolled
		c.JSON(http.StatusConflict, gin.H{
			"error": "Student is already enrolled in this course or invalid IDs provided",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Enrollment successful"})
}

// GetByUser fetches all courses a specific student is enrolled in
func (h *EnrollmentHandler) GetByUser(c *gin.Context) {
	userID := c.Param("user_id")

	query := `
		SELECT e.id, c.id, c.title, e.enrolled_at 
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE e.user_id = $1`

	rows, err := h.DB.QueryContext(c.Request.Context(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var eid, cid int
		var title string
		var enrolledAt interface{}
		rows.Scan(&eid, &cid, &title, &enrolledAt)
		results = append(results, map[string]interface{}{
			"enrollment_id": eid,
			"course_id":     cid,
			"course_title":  title,
			"enrolled_at":   enrolledAt,
		})
	}

	c.JSON(http.StatusOK, results)
}

// GetByCourse fetches all students enrolled in a specific course
func (h *EnrollmentHandler) GetByCourse(c *gin.Context) {
	courseID := c.Param("course_id")

	query := `
		SELECT e.id, u.id, u.first_name, u.last_name, e.enrolled_at 
		FROM enrollments e
		JOIN users u ON e.user_id = u.id
		WHERE e.course_id = $1`

	rows, err := h.DB.QueryContext(c.Request.Context(), query, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course roster"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var eid, uid int
		var fname, lname string
		var enrolledAt interface{}
		rows.Scan(&eid, &uid, &fname, &lname, &enrolledAt)
		results = append(results, map[string]interface{}{
			"enrollment_id": eid,
			"user_id":       uid,
			"full_name":     fname + " " + lname,
			"enrolled_at":   enrolledAt,
		})
	}

	c.JSON(http.StatusOK, results)
}

// Unenroll handles removing a student from a course
func (h *EnrollmentHandler) Unenroll(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM enrollments WHERE id = $1`
	affected, err := h.DB.Exec(c.Request.Context(), query, id)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unenroll"})
		return
	}

	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unenrolled"})
}