package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user entities.User
	initializers.DB.Where("id=?", claims["id"]).Find(&user)

	if user.ID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser", claims["user"])
	c.Next()
}
