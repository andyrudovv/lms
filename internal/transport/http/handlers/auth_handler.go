package handlers

import (
	"net/http"

	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/dto"
	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.auth.RegisterStudent(c.Request.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	responder.Created(c, gin.H{"id": id})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responder.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		responder.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	responder.OK(c, gin.H{"access_token": token})
}
