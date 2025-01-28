package router

import (
	"github.com/bartholomeas/hwheels_api/config/initializers"
	controllers "github.com/bartholomeas/hwheels_api/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	authController := controllers.NewAuthController(initializers.DB)

	auth.POST("/register", authController.CreateUser)
	auth.POST("/login", authController.LoginUser)
}
