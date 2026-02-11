package middleware

import (
	"net/http"
	"strings"

	"lms-backend/internal/service"
	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"
const CtxRoleKey = "role"

func AuthJWT(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			responder.Fail(c, http.StatusUnauthorized, "missing bearer token")
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		claims, err := auth.Parse(tokenStr)
		if err != nil {
			responder.Fail(c, http.StatusUnauthorized, "invalid token")
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)
		c.Next()
	}
}
