package handlers

import (
	"net/http"
	"strconv"
	"time"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/dto"
	"lms-backend/internal/transport/http/middleware"
	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.svc.Create(c.Request.Context(), model.User{
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}, req.Password)
	if err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.Created(c, gin.H{"id": id})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]gin.H, 0, len(users))
	for _, u := range users {
		out = append(out, gin.H{"id": u.ID, "email": u.Email, "full_name": u.FullName, "role_id": u.RoleID})
	}

	responder.OK(c, gin.H{"items": out, "count": len(out), "ts": time.Now()})
}

// Admin: change role of user by id (body: { "role": "teacher" })
func (h *UserHandler) ChangeRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		responder.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.ChangeRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.ChangeRole(c.Request.Context(), userID, req.Role); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.OK(c, gin.H{"status": "updated"})
}

// Any logged-in user: returns own profile based on JWT
func (h *UserHandler) Me(c *gin.Context) {
	uidAny, _ := c.Get(middleware.CtxUserIDKey)
	uid, _ := uidAny.(int)

	u, role, err := h.svc.Me(c.Request.Context(), uid)
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	responder.OK(c, gin.H{
		"id": u.ID,
		"email": u.Email,
		"full_name": u.FullName,
		"role": role,
		"role_id": u.RoleID,
	})
}

func (h *UserHandler) Roles(c *gin.Context) {
	roles, err := h.svc.ListRoles(c.Request.Context())
	if err != nil {
		responder.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]gin.H, 0, len(roles))
	for _, r := range roles {
		out = append(out, gin.H{"id": r.ID, "name": r.Name})
	}
	responder.OK(c, gin.H{"items": out, "count": len(out)})
}
