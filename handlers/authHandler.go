package handlers

import (
	"gitrankhub/config"
	"gitrankhub/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterationInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func createToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// create token now
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString([]byte(secretKey))

	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func VerifyToken(tokenString string) (bool, error) {
	// return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	// validate the alg is what you expect
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, jwt.ErrSignatureInvalid
	// 	}
	// 	return []byte(secretKey), nil
	// })
	token, err := ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	return token.Valid, err
}

// extractTokenUserID extracts the user ID from the JWT token claims
func ExtractTokenUserID(tokenString string) (string, error) {
	accessToken, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok || !accessToken.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	return userID, nil
}

// LoginUser handles user login and JWT token generation
func LoginUser(c *gin.Context) {
	var loginInput LoginInput
	if err := c.BindJSON(&loginInput); err != nil {
		log.Println("Error binding JSON:", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user, err := gorm.G[models.User](config.ConnectDB()).
		Where("email = ?", loginInput.Email).
		First(ctx)

	if err != nil {
		log.Println("Error finding user:", err)
		c.IndentedJSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := createToken(strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		log.Println("Error creating token:", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to create access token"})
		return
	}

	if !VerifyPassword(loginInput.Password, user.Password) {
		log.Println("Invalid password for user:", loginInput.Email)
		c.IndentedJSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "access_token": accessToken})
}

// RegisterUser handles new user registration
func RegisterUser(c *gin.Context) {
	var newUser models.User
	var regInput RegisterationInput

	if err := c.BindJSON(&regInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser.Username = regInput.Username
	newUser.Email = regInput.Email
	passwordHash, err := HashPassword(regInput.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = passwordHash

	err = gorm.G[models.User](config.ConnectDB()).Create(ctx, &newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newUser)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
