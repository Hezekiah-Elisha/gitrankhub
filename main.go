package main

import (
	"gitrankhub/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users", handlers.GetUsers)
	router.POST("/users", handlers.CreateUser)
	router.GET("/users/:id", handlers.GetUserByID)
	router.GET("/users/email/:email", handlers.GetUserByEmail)
	router.PUT("/users/:id", handlers.UpdateUser)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.Run(":8080")
}
