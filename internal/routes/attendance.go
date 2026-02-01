package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Attendance(rg *gin.RouterGroup) {
	attendance := rg.Group("/attendance")
	{
		attendance.POST("", markAttendance)
		attendance.GET("/course/:course_id", getAttendanceByCourse)
		attendance.GET("/student/:student_id", getAttendanceByStudent)
	}
}

func markAttendance(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "attendance marked",
	})
}

func getAttendanceByCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"course_id":  c.Param("course_id"),
		"attendance": []string{},
	})
}

func getAttendanceByStudent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"student_id": c.Param("student_id"),
		"attendance": []string{},
	})
}
