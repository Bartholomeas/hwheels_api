package router

import (
	"github.com/bartholomeas/hwheels_api/internal/user/controller"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	userController := controller.NewUserController()

	// user.GET("/profile", middlewares.CheckAuth, userController.GetUserProfile)
	user.GET("/profile", userController.GetUserByToken)
}
