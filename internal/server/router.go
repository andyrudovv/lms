package server

import (
	"github.com/gin-gonic/gin"
	"github.com/andyrudovv/lms/internal/routes"
	"github.com/andyrudovv/lms/internal/infrastructure/db"
)

func RegisterRoutes(r *gin.Engine, database *db.PostgresConnector) {
	routes.Health(r)

	api := r.Group("/api/v1")
	{
		routes.Auth(api, database)
		routes.Courses(api, database)
		routes.Attendance(api, database)
		routes.Users(api, database)
		routes.Submissions(api, database)
	}
}