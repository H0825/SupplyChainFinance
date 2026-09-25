package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRoles validates role from X-User-Role header.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allow := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allow[strings.ToLower(strings.TrimSpace(r))] = struct{}{}
	}

	return func(c *gin.Context) {
		role := strings.ToLower(strings.TrimSpace(c.GetHeader("X-User-Role")))
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing role header: X-User-Role"})
			return
		}
		if _, ok := allow[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not allowed"})
			return
		}
		c.Set("role", role)
		c.Next()
	}
}
