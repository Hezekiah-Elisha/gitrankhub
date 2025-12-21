package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gitrankhub/config"
	"gitrankhub/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")

type Result struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterationInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GitUserDetails struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`

	Name            *string `json:"name"`
	Company         *string `json:"company"`
	Blog            *string `json:"blog"`
	Location        *string `json:"location"`
	Email           *string `json:"email"`
	Hireable        *bool   `json:"hireable"`
	Bio             *string `json:"bio"`
	TwitterUsername *string `json:"twitter_username"`

	PublicRepos int    `json:"public_repos"`
	PublicGists int    `json:"public_gists"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
		return []byte(secretKey), nil
	})
}

func VerifyToken(tokenString string) (bool, error) {
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
		First(context.Background())

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

	err = gorm.G[models.User](config.ConnectDB()).Create(context.Background(), &newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newUser)
}

// GithubCallBack handles GitHub OAuth callback
func GithubCallBack(c *gin.Context) {
	code := c.Query("code")

	resp, err := fetchGithubAccessToken(code)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access token"})
		return
	}

	userDetails, err := getUserDetailsFromGithub(resp)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	user, err := gorm.G[models.User](config.ConnectDB()).Where("username = ?", userDetails.Login).First(context.Background())
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Error finding user:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	if code == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": user, "code": code})
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with its hashed version
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// fetchGithubAccessToken exchanges the authorization code for an access token
func fetchGithubAccessToken(code string) (string, error) {
	var result Result

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("missing GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch access token, status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("no access token found in response")
	}

	return result.AccessToken, nil
}

func getUserDetailsFromGithub(accessToken string) (GitUserDetails, error) {
	var userDetails GitUserDetails

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return userDetails, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "git-rank-hub")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return userDetails, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		return userDetails, fmt.Errorf(
			"github api error: %d - %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&userDetails); err != nil {
		return userDetails, err
	}

	log.Println("GitHub User Details:", userDetails) // Log the actual struct here

	return userDetails, nil
}
