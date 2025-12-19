package middlewares

import (
	"gitrankhub/handlers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		// Strip "Bearer " prefix
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		valid, err := handlers.VerifyToken(tokenString)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID, err := handlers.ExtractTokenUserID(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not extract user ID from token"})
			return
		}
		c.Set("userID", userID)

		c.Next()
	}
}
