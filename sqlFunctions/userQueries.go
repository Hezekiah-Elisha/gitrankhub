package sqlfunctions

import (
	"context"
	"gitrankhub/config"
	"gitrankhub/models"

	"gorm.io/gorm"
)

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	user, err := gorm.G[models.User](config.ConnectDB()).
		Where("email = ?", email).
		First(context.Background())

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserByID(userID string) (models.User, error) {
	user, err := gorm.G[models.User](config.ConnectDB()).
		Where("id = ?", userID).
		First(context.Background())

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	result := config.ConnectDB().Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func CreateUser(username, name, email, password, role, avatarURL, bio string) error {
	// Placeholder function - implement actual SQL insert logic here
	return nil
}

func UpdateUser(userID, username, name, email, password, role, avatarURL, bio string) error {
	// Placeholder function - implement actual SQL update logic here
	return nil
}

func DeleteUser(userID string) error {
	// Placeholder function - implement actual SQL delete logic here
	return nil
}
