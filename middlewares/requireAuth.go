package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"rakamin.com/final-task/database"
	"rakamin.com/final-task/models"
)

func RequireAuth(c *gin.Context) {
	// get the cookie
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// decode/ validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check exp token
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find user with token
		var user models.User
		database.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// attach to req
		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
