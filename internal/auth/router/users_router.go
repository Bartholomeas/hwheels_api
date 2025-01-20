package router

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/internal/auth/controllers"
	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
}

func InitAuthRouter(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	auth.POST("/register", controllers.CreateUser)
	auth.POST("/login", LoginUser)
}
