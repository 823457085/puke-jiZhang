package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简化版认证：前端传入 user_id（生产环境需用真实微信session）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// MVP阶段：从header获取userID
		// 生产环境应使用微信授权后的session token
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID header"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
