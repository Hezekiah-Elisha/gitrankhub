package handlers

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

type LoginInput struct {
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

func verifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
}

func extractTokenUserID(accessToken *jwt.Token) (string, error) {
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

func LoginUser(c *gin.Context) {
	var loginInput LoginInput
	if err := c.BindJSON(&loginInput); err != nil {
		log.Println("Error binding JSON:", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid input"})
		return
	}
	//return the input for testing purpose
	c.IndentedJSON(200, gin.H{"message": "Login successful", "email": loginInput.Email, "password": loginInput.Password})
}
