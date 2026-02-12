package handlers

import (
	"log"
	"net/http"
	"strconv"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/dto"
	"lms-backend/internal/transport/http/middleware"
	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(svc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

func (h *AttendanceHandler) Mark(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	var req dto.MarkAttendanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.svc.Mark(c.Request.Context(), model.Attendance{
		CourseID:   courseID,
		StudentID:  req.StudentID,
		LessonDate: req.LessonDate,
		Status:     req.Status,
		Note:       req.Note,
	})
	if err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.OK(c, gin.H{"status": "saved"})
}

func (h *AttendanceHandler) ListByCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	items, err := h.svc.ListByCourse(c.Request.Context(), courseID)
	if err != nil {
		log.Println("Error listing attendance by course:", err)
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, a := range items {
		out = append(out, gin.H{
			"id": a.ID, "course_id": a.CourseID, "student_id": a.StudentID,
			"lesson_date": a.LessonDate, "status": a.Status, "note": a.Note,
		})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out)})
}

// Student: my attendance (optional ?course_id=)
func (h *AttendanceHandler) MyAttendance(c *gin.Context) {
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	uid, _ := uidAny.(int)
	roleAny, _ := c.Get(middleware.CtxRoleKey)
	role, _ := roleAny.(string)

	if role != "student" && role != "admin" && role != "teacher" {
		responder.Fail(c, http.StatusForbidden, "forbidden")
		return
	}

	courseID := 0
	if v := c.Query("course_id"); v != "" {
		x, err := strconv.Atoi(v)
		if err != nil || x < 0 {
			responder.Fail(c, http.StatusBadRequest, "invalid course_id")
			return
		}
		courseID = x
	}

	// If admin/teacher calls this, it returns their own (usually empty). For FE, intended for students.
	items, err := h.svc.ListByStudent(c.Request.Context(), uid, courseID)
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, a := range items {
		out = append(out, gin.H{
			"id": a.ID, "course_id": a.CourseID, "student_id": a.StudentID,
			"lesson_date": a.LessonDate, "status": a.Status, "note": a.Note,
		})
	}
	responder.OK(c, gin.H{"items": out, "count": len(out)})
}
