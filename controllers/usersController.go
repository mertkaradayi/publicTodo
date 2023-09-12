package controllers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"examples.com/jwt-auth/initializers"
	"examples.com/jwt-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get the request body
	var body struct {
		Email    string
		Password string
		Name     string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate email format
	if !isValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	// Validate password complexity
	if !isValidPassword(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password format",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
		Name:     body.Name,
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	//Get the email and password of req body
	var body struct {
		Email    string
		Password string
		Name     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}

	//Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//Generate a jwt token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	fmt.Println(tokenString, err)

	if err != nil {
		// Handle error
		fmt.Println("Error generating token:", err)
		return
	}

	//Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func Logout(c *gin.Context) {
	// Get the token from the cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Add the token to the token_blacklist table with an expiration time
	// This will prevent the token from being valid after logging out
	expirationTime := time.Now().Add(24 * time.Hour) // Set an appropriate expiration time
	invalidatedToken := models.BlackListedToken{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}
	result := initializers.DB.Create(&invalidatedToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to invalidate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func GetUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Create a response payload containing user's information
	userProfile := struct {
		ID    uint
		Email string
		Name  string
	}{
		ID:    user.(models.User).ID,
		Email: user.(models.User).Email,
		Name:  user.(models.User).Name,
	}

	c.JSON(http.StatusOK, userProfile)
}

func DeleteUser(c *gin.Context) {
	// Get the user from the context (set by RequireAuth middleware)
	user, _ := c.Get("user")

	// Get the user's ID from the context
	userID := user.(models.User).ID

	// Delete the user from the database
	result := initializers.DB.Delete(&models.User{}, userID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func ChangePassword(c *gin.Context) {
	// Get the user from the context (set by RequireAuth middleware)
	user, _ := c.Get("user")

	// Get the user's ID from the context
	userID := user.(models.User).ID

	// Bind the new password from the request body
	var requestBody struct {
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Validate the new password's complexity
	if !isValidPassword(requestBody.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password format",
		})
		return
	}

	// Hash the new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash new password",
		})
		return
	}

	// Update the user's password in the database
	result := initializers.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", string(newHash))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update password",
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

func isValidEmail(email string) bool {
	// Use a regular expression to validate email format
	// You can adjust this regex pattern to your needs
	// This is a simple example, more comprehensive patterns are available
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

func isValidPassword(password string) bool {
	// Add password complexity checks here (e.g., length, special characters)
	// For example, check if the password is at least 8 characters long
	return len(password) >= 8
}
