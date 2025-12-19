package middlewares

import (
	"gitrankhub/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		valid, err := handlers.VerifyToken(tokenString)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// userID, err := handlers.ExtractTokenUserID(tokenString)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not extract user ID from token"})
		// 	return
		// }
		// c.Set("userID", userID)

		c.Next()
	}
}
