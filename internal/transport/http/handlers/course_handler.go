package handlers

import (
	"net/http"
	"strconv"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/dto"
	"lms-backend/internal/transport/http/middleware"
	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	svc *service.CourseService
}

func NewCourseHandler(svc *service.CourseService) *CourseHandler {
	return &CourseHandler{svc: svc}
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CreateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.TeacherID == 0 {
		uidAny, _ := c.Get(middleware.CtxUserIDKey)
		req.TeacherID, _ = uidAny.(int)
	}

	id, err := h.svc.Create(c.Request.Context(), model.Course{Title: req.Title, TeacherID: req.TeacherID})
	if err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.Created(c, gin.H{"id": id})
}

func (h *CourseHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, x := range items {
		out = append(out, gin.H{"id": x.ID, "title": x.Title, "teacher_id": x.TeacherID, "created_at": x.CreatedAt})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out)})
}

func (h *CourseHandler) Enroll(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	var req dto.EnrollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Enroll(c.Request.Context(), courseID, req.StudentID); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.OK(c, gin.H{"status": "enrolled"})
}

// Any logged-in user: returns my courses based on role
func (h *CourseHandler) MyCourses(c *gin.Context) {
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	uid, _ := uidAny.(int)
	roleAny, _ := c.Get(middleware.CtxRoleKey)
	role, _ := roleAny.(string)

	var items []model.Course
	var err error

	switch role {
	case "student":
		items, err = h.svc.ListByStudent(c.Request.Context(), uid)
	case "teacher":
		items, err = h.svc.ListByTeacher(c.Request.Context(), uid)
	case "admin":
		items, err = h.svc.List(c.Request.Context())
	default:
		responder.Fail(c, http.StatusForbidden, "unknown role")
		return
	}
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, x := range items {
		out = append(out, gin.H{"id": x.ID, "title": x.Title, "teacher_id": x.TeacherID, "created_at": x.CreatedAt})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out)})
}


func (h *CourseHandler) Unenroll(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	var req dto.EnrollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Unenroll(c.Request.Context(), courseID, req.StudentID); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.OK(c, gin.H{"status": "unenrolled"})
}

func (h *CourseHandler) GetStudents(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	items, err := h.svc.GetStudents(c.Request.Context(), courseID)
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, s := range items {
		out = append(out, gin.H{"id": s.ID, "full_name": s.FullName, "email": s.Email})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out)})
}

func (h *CourseHandler) GetAvailableStudents(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil || courseID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid course id")
		return
	}

	items, err := h.svc.GetAvailableStudents(c.Request.Context(), courseID)
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, s := range items {
		out = append(out, gin.H{"id": s.ID, "full_name": s.FullName, "email": s.Email})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out)})
}
