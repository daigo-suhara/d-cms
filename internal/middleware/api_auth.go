package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIKeyValidator interface {
	Validate(ctx context.Context, key string) bool
}

// APIAuth requires a valid API key via "Authorization: Bearer <key>" header.
func APIAuth(validator APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "APIキーが必要です。Authorization: Bearer <key> ヘッダーを付与してください"})
			return
		}
		key := strings.TrimPrefix(authHeader, "Bearer ")
		if !validator.Validate(c.Request.Context(), key) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なAPIキーです"})
			return
		}
		c.Next()
	}
}
