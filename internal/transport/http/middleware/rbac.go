package middleware

import (
	"net/http"

	"lms-backend/internal/transport/http/responder"

	"github.com/gin-gonic/gin"
)

func RequireRoles(allowed ...string) gin.HandlerFunc {
	set := map[string]struct{}{}
	for _, r := range allowed {
		set[r] = struct{}{}
	}

	return func(c *gin.Context) {
		roleAny, ok := c.Get(CtxRoleKey)
		if !ok {
			responder.Fail(c, http.StatusForbidden, "no role in context")
			return
		}
		role, _ := roleAny.(string)

		if _, ok := set[role]; !ok {
			responder.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		c.Next()
	}
}
