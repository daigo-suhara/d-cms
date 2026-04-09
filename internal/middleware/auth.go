package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth validates requests using a static Bearer token.
// HTML requests are redirected to /admin/login on failure; JSON requests receive 401.
func Auth(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			if strings.TrimPrefix(authHeader, "Bearer ") == token {
				c.Next()
				return
			}
		}

		// Check session cookie (for browser-based admin UI)
		cookie, err := c.Cookie("admin_token")
		if err == nil && cookie == token {
			c.Next()
			return
		}

		accept := c.GetHeader("Accept")
		if strings.Contains(accept, "text/html") {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}
