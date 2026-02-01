package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/andyrudovv/lms/internal/infrastructure/db" // Adjust path as needed
)

type AuthHandler struct {
	DB *db.PostgresConnector
}

func Auth(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &AuthHandler{DB: database}
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/register", h.register)
	}
}

func (h *AuthHandler) register(c *gin.Context) {
	var input struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		RoleID    int    `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users (first_name, last_name, email, password_hash, role_id) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	var id int
	err := h.DB.QueryRowContext(c.Request.Context(), query, 
		input.FirstName, input.LastName, input.Email, input.Password, input.RoleID).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "user registered"})
}

func (h *AuthHandler) login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedHash string
	query := `SELECT password_hash FROM users WHERE email = $1`
	err := h.DB.QueryRowContext(c.Request.Context(), query, input.Email).Scan(&storedHash)

	if err != nil || storedHash != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}