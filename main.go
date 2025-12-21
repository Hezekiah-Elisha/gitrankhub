package main

import (
	"gitrankhub/handlers"
	"gitrankhub/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS middleware must be registered before routes
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "OPTIONS", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handlers.LoginUser)
		authRoutes.POST("/register", handlers.RegisterUser)
		authRoutes.GET("/github/callback", handlers.GithubCallBack)
	}
	userRoutes := router.Group("/users")
	{
		// Protected routes can be added here with middleware if needed
		userRoutes.GET("/", middlewares.AuthMiddleware(), handlers.GetUsers)
		userRoutes.POST("/", middlewares.AuthMiddleware(), handlers.CreateUser)
		userRoutes.GET("/:id", middlewares.AuthMiddleware(), handlers.GetUserByID)
		userRoutes.GET("/email/:email", middlewares.AuthMiddleware(), handlers.GetUserByEmail)
		userRoutes.PUT("/:id", middlewares.AuthMiddleware(), handlers.UpdateUser)
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
		})
	})

	router.Run(":8080")
}
