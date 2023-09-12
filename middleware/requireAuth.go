package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"examples.com/jwt-auth/initializers"
	"examples.com/jwt-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {

	//get the cookie of req
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	// sample token string taken from the New example

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)

		}
		//Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)

		}

		// Check if the token exists in the token_blacklist table
		var blacklistedToken models.BlackListedToken
		result := initializers.DB.Where("token = ?", tokenString).First(&blacklistedToken)
		if result.Error == nil && time.Now().Before(blacklistedToken.ExpiresAt) {
			// Token is invalidated, deny the request
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//Attach to req
		c.Set("user", user)
		//Continue
		fmt.Println("In middleware")
		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
