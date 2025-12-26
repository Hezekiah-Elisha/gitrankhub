package handlers

import (
	"context"
	"gitrankhub/config"
	"gitrankhub/models"
	sqlfunctions "gitrankhub/sqlFunctions"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateUserInput struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

var ctx = context.Background()

func GetUsers(c *gin.Context) {
	var users []models.User

	users, err := sqlfunctions.GetAllUsers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(200, users)
}

func CreateUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := gorm.G[models.User](config.ConnectDB()).Create(ctx, &newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := sqlfunctions.GetUserByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := sqlfunctions.GetUserByEmail(email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedData UpdateUserInput

	if err := c.BindJSON(&updatedData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := gorm.G[models.User](config.ConnectDB()).
		Where("id = ?", id).
		First(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Build a dynamic update payload only with provided fields
	if updatedData.Username != nil {
		user.Username = *updatedData.Username
	}
	if updatedData.Email != nil {
		user.Email = *updatedData.Email
	}
	if updatedData.Password != nil {
		user.Password = *updatedData.Password
	}

	_, err = gorm.G[models.User](config.ConnectDB()).Where("id = ?", user.ID).Updates(ctx, user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch and return the full updated user object
	updatedUser, err := gorm.G[models.User](config.ConnectDB()).Where("id = ?", user.ID).First(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}
