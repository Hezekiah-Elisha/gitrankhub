package middlewares

import "github.com/gin-gonic/gin"

// router.Use(cors.New(cors.Config{
// 	AllowOrigins:     []string{"http://localhost:3000"},
// 	AllowMethods:     []string{"PUT", "OPTIONS", "PATCH", "GET", "POST", "DELETE"},
// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 	ExposeHeaders:    []string{"Content-Length"},
// 	AllowCredentials: true,
// 	MaxAge:           12 * time.Hour,
// }))
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}

}
