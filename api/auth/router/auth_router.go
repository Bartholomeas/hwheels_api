package router

import (
	controllers "github.com/bartholomeas/hwheels_api/api/auth/controller"
	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	authController := controllers.NewAuthController(initializers.DB)

	auth.POST("/register", authController.CreateUser)
	auth.POST("/login", authController.LoginUser)
}
