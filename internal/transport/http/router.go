package httpapi

import (
	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/handlers"
	"lms-backend/internal/transport/http/middleware"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authSvc *service.AuthService,
	authH *handlers.AuthHandler,
	userH *handlers.UserHandler,
	courseH *handlers.CourseHandler,
	attH *handlers.AttendanceHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestLogger(), gin.Recovery(), middleware.ErrorHandler())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")

	// public
	api.POST("/auth/register", authH.Register) // creates student
	api.POST("/auth/login", authH.Login)

	// protected
	protected := api.Group("/")
	protected.Use(middleware.AuthJWT(authSvc))
	{
		// profile
		protected.GET("/me", userH.Me)
		protected.GET("/roles", middleware.RequireRoles("admin"), userH.Roles)

		// users (admin)
		protected.POST("/users", middleware.RequireRoles("admin"), userH.Create)
		protected.GET("/users", middleware.RequireRoles("admin"), userH.List)
		protected.PATCH("/users/:id/role", middleware.RequireRoles("admin"), userH.ChangeRole)

		// courses
		protected.POST("/courses", middleware.RequireRoles("admin", "teacher"), courseH.Create)
		protected.GET("/courses", middleware.RequireRoles("admin", "teacher", "student"), courseH.List)
		protected.GET("/my/courses", middleware.RequireRoles("admin", "teacher", "student"), courseH.MyCourses)
		protected.POST("/courses/:id/enroll", middleware.RequireRoles("admin", "teacher"), courseH.Enroll)

		// attendance
		protected.POST("/courses/:id/attendance", middleware.RequireRoles("admin", "teacher"), attH.Mark)
		protected.GET("/courses/:id/attendance", middleware.RequireRoles("admin", "teacher"), attH.ListByCourse)
		protected.GET("/my/attendance", middleware.RequireRoles("admin", "teacher", "student"), attH.MyAttendance)
	}

	return r
}
