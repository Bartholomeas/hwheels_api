package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
}

func InitUsersRouter(router *gin.RouterGroup) {
	user := router.Group("/auth")

	user.POST("/register", RegisterUser)
	user.POST("/login", LoginUser)
}
