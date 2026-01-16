package server

import (
	"github.com/gin-gonic/gin"
	"github.com/andyrudovv/lms/internal/routes"
)

func RegisterRoutes(r *gin.Engine) {
	// Public routes
	routes.Health(r)

	api := r.Group("/api/v1")
	{
		routes.Auth(api)
		routes.Courses(api)
	}
}
