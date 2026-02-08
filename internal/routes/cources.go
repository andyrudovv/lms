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

// Courses registers all course, module, enrollment, and progress routes
func Courses(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &CourseHandler{DB: database}
	
	courses := rg.Group("/courses")
	{
		// Course Management
		courses.GET("", h.GetCourses)
		courses.POST("", h.CreateCourse)
		courses.GET("/:id", h.GetCourseByID)

		// Structure (Modules/Weeks & Lessons/Resources)
		courses.POST("/modules", h.CreateStructure)
		courses.POST("/modules/:week_id/lessons", h.AddLesson)

		// Enrollment Logic
		courses.POST("/:id/enroll", h.EnrollStudent)

		// Progress Tracking & Calculations
		courses.GET("/:id/progress", h.GetProgress)
	}
}

// ==========================
// COURSE MANAGEMENT
// ==========================

func (h *CourseHandler) CreateCourse(c *gin.Context) {
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

	// Responsibility: Async Notification/Publishing Log
	go func(courseTitle string) {
		time.Sleep(2 * time.Second)
		println("Background Task: Course published and notifications sent for:", courseTitle)
	}(newCourse.Title)

	c.JSON(http.StatusCreated, newCourse)
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
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

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	var crs models.Course
	query := `SELECT id, title, description, teacher_id, created_at FROM courses WHERE id = $1`
	
	err := h.DB.QueryRowContext(c.Request.Context(), query, id).Scan(
		&crs.ID, &crs.Title, &crs.Description, &crs.TeacherID, &crs.CreatedAt)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, crs)
}

// ==========================
// COURSE STRUCTURE (MODULES & LESSONS)
// ==========================

func (h *CourseHandler) CreateStructure(c *gin.Context) {
	var input struct {
		CourseID int    `json:"course_id" binding:"required"`
		Title    string `json:"title" binding:"required"`
		Number   int    `json:"week_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO weeks (course_id, week_number, title) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := h.DB.QueryRowContext(c.Request.Context(), query, input.CourseID, input.Number, input.Title).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create module"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "status": "module created"})
}

func (h *CourseHandler) AddLesson(c *gin.Context) {
	weekID := c.Param("week_id")
	var input struct {
		Title string `json:"title" binding:"required"`
		Type  string `json:"type" binding:"required"` // e.g., video, pdf, link
		URL   string `json:"url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO resources (week_id, title, type, url) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := h.DB.QueryRowContext(c.Request.Context(), query, weekID, input.Title, input.Type, input.URL).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add lesson content"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "status": "lesson added to module"})
}

// ==========================
// ENROLLMENT & PROGRESS
// ==========================

func (h *CourseHandler) EnrollStudent(c *gin.Context) {
	courseID := c.Param("id")
	var input struct {
		UserID int `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO enrollments (user_id, course_id) VALUES ($1, $2)`
	_, err := h.DB.Exec(c.Request.Context(), query, input.UserID, courseID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Student already enrolled or invalid data"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Enrolled successfully"})
}

func (h *CourseHandler) GetProgress(c *gin.Context) {
	courseID := c.Param("id")
	studentID := c.Query("student_id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required as query param"})
		return
	}

	// Calculation: Ratio of attended modules to total modules in the course
	query := `
		SELECT 
			(SELECT COUNT(*) FROM weeks WHERE course_id = $1) as total_modules,
			(SELECT COUNT(DISTINCT date) FROM attendance 
			 WHERE course_id = $1 AND student_id = $2 AND status = 'present') as modules_completed
	`
	
	var total, completed int
	err := h.DB.QueryRowContext(c.Request.Context(), query, courseID, studentID).Scan(&total, &completed)
	
	if err != nil || total == 0 {
		c.JSON(http.StatusOK, gin.H{"progress_percentage": 0, "message": "No modules found for this course"})
		return
	}

	percentage := (float64(completed) / float64(total)) * 100
	c.JSON(http.StatusOK, gin.H{
		"course_id":           courseID,
		"student_id":          studentID,
		"total_modules":       total,
		"completed_modules":   completed,
		"progress_percentage": percentage,
	})
}