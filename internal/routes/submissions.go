package routes

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	DB *db.PostgresConnector
}

func Submissions(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &SubmissionHandler{DB: database}
	sub := rg.Group("/submissions")
	{
		sub.POST("/upload", h.submitExam)
		sub.GET("/exam/:exam_id", h.getByExam)
		sub.PATCH("/grade/:id", h.gradeSubmission) // To update 'grade' and 'feedback'
	}
}

func (h *SubmissionHandler) submitExam(c *gin.Context) {
	examID := c.PostForm("exam_id")
	studentID := c.PostForm("student_id")

	// Handle File Upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Save file to a 'uploads' directory
	path := filepath.Join("uploads", fmt.Sprintf("exam_%s_student_%s_%s", examID, studentID, file.Filename))
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Insert into DB using ON CONFLICT as per your unique index idx_submissions_exam_student
	query := `INSERT INTO submissions (exam_id, student_id, file_path) 
	          VALUES ($1, $2, $3) 
	          ON CONFLICT (exam_id, student_id) 
	          DO UPDATE SET file_path = EXCLUDED.file_path, submitted_at = now()`
	
	_, err = h.DB.Exec(c.Request.Context(), query, examID, studentID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Exam submitted successfully", "path": path})
}

func (h *SubmissionHandler) gradeSubmission(c *gin.Context) {
	submissionID := c.Param("id")

	var input struct {
		Grade    int    `json:"grade" binding:"required"`
		Feedback string `json:"feedback"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the grade and feedback for the specific submission ID
	query := `UPDATE submissions 
	          SET grade = $1, feedback = $2 
	          WHERE id = $3`
	
	result, err := h.DB.Exec(c.Request.Context(), query, input.Grade, input.Feedback, submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grade"})
		return
	}

	if result == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission graded successfully"})
}


func (h *SubmissionHandler) getByExam(c *gin.Context) {
	examID := c.Param("exam_id")

	query := `SELECT id, exam_id, student_id, file_path, submitted_at, grade, feedback 
	          FROM submissions WHERE exam_id = $1`
	
	rows, err := h.DB.QueryContext(c.Request.Context(), query, examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, eID, sID int
		var filePath, feedback *string
		var grade *int
		var submittedAt interface{}
		
		rows.Scan(&id, &eID, &sID, &filePath, &submittedAt, &grade, &feedback)
		results = append(results, map[string]interface{}{
			"id":           id,
			"student_id":   sID,
			"file_path":    filePath,
			"submitted_at": submittedAt,
			"grade":        grade,
			"feedback":     feedback,
		})
	}

	c.JSON(http.StatusOK, results)
}