package routes

import (
	"net/http"

	"github.com/andyrudovv/lms/internal/infrastructure/db"
	"github.com/andyrudovv/lms/internal/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *db.PostgresConnector
}

// Users registers the user-related routes and initializes the handler
func Users(rg *gin.RouterGroup, database *db.PostgresConnector) {
	h := &UserHandler{DB: database}

	users := rg.Group("/users")
	{
		users.GET("", h.getUsers)
		users.GET("/:id", h.getUserByID)
		users.PUT("/:id", h.updateUser)
	}
}

// getUsers fetches all users with their role names joined from the roles table
func (h *UserHandler) getUsers(c *gin.Context) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.role_id, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.id ASC`

	rows, err := h.DB.QueryContext(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var u models.User
		var roleName string
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.RoleID, &roleName)
		if err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":         u.ID,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"email":      u.Email,
			"role_id":    u.RoleID,
			"role_name":  roleName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// getUserByID fetches a single user by their primary key
func (h *UserHandler) getUserByID(c *gin.Context) {
	id := c.Param("id")

	var u models.User
	query := `SELECT id, first_name, last_name, email, role_id FROM users WHERE id = $1`

	err := h.DB.QueryRowContext(c.Request.Context(), query, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.RoleID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// updateUser allows updating user profile information
func (h *UserHandler) updateUser(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, updated_at = now() 
		WHERE id = $4`

	rowsAffected, err := h.DB.Exec(c.Request.Context(), query, input.FirstName, input.LastName, input.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}