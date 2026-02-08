package routes

import (
	"net/http"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	DB *db.PostgresConnector
}

func Content(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &ContentHandler{DB: database}

	// Modules
	rg.POST("/courses/:id/modules", h.createModule)
	rg.GET("/courses/:id/modules", h.getModules)

	// Lessons
	rg.POST("/modules/:week_id/lessons", h.createLesson)
	rg.GET("/modules/:week_id/lessons", h.getLessons)
}

func (h *ContentHandler) createModule(c *gin.Context) {
	courseID := c.Param("id")
	var input struct {
		WeekNumber int    `json:"week_number" binding:"required"`
		Title      string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `INSERT INTO weeks (course_id, week_number, title) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := h.DB.QueryRowContext(c.Request.Context(), query, courseID, input.WeekNumber, input.Title).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Module created"})
}

func (h *ContentHandler) createLesson(c *gin.Context) {
	weekID := c.Param("week_id")
	var input struct {
		Title string `json:"title" binding:"required"`
		Type  string `json:"type" binding:"required"`
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Lesson added"})
}

func (h *ContentHandler) getModules(c *gin.Context) {
	courseID := c.Param("id")
	query := `SELECT id, week_number, title FROM weeks WHERE course_id = $1 ORDER BY week_number ASC`
	rows, err := h.DB.QueryContext(c.Request.Context(), query, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var modules []map[string]interface{}
	for rows.Next() {
		var id, weekNumber int
		var title string
		err := rows.Scan(&id, &weekNumber, &title)
		if err != nil {
			continue
		}
		modules = append(modules, map[string]interface{}{
			"id":          id,
			"week_number": weekNumber,
			"title":       title,
		})
	}
	c.JSON(http.StatusOK, gin.H{"modules": modules})
}

func (h *ContentHandler) getLessons(c *gin.Context) {
	weekID := c.Param("week_id")
	query := `SELECT id, title, type, url FROM resources WHERE week_id = $1 ORDER BY id ASC`
	rows, err := h.DB.QueryContext(c.Request.Context(), query, weekID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var lessons []map[string]interface{}
	for rows.Next() {
		var id int
		var title string
		var lessonType string
		var url string
		err := rows.Scan(&id, &title, &lessonType, &url)
		if err != nil {
			continue
		}
		lessons = append(lessons, map[string]interface{}{
			"id":    id,
			"title": title,
			"type":  lessonType,
			"url":   url,
		})
	}
	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}