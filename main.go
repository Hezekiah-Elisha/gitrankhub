package main

import (
	"gitrankhub/handlers"
	"gitrankhub/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handlers.LoginUser)
		authRoutes.POST("/register", handlers.RegisterUser)
	}
	userRoutes := router.Group("/users")
	{
		// Protected routes can be added here with middleware if needed
		userRoutes.GET("/", middlewares.AuthMiddleware(), handlers.GetUsers)
		userRoutes.POST("/", handlers.CreateUser)
		userRoutes.GET("/:id", handlers.GetUserByID)
		userRoutes.GET("/email/:email", handlers.GetUserByEmail)
		userRoutes.PUT("/:id", handlers.UpdateUser)
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
		})
	})

	router.Run(":8080")
}
