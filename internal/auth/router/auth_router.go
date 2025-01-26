package router

import (
	controllers "github.com/bartholomeas/hwheels_api/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	auth.POST("/register", controllers.CreateUser)
	auth.POST("/login", controllers.LoginUser)
}
