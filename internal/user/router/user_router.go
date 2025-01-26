package router

import (
	"github.com/bartholomeas/hwheels_api/internal/common/middlewares"
	"github.com/bartholomeas/hwheels_api/internal/user/controller"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")

	user.GET("/profile", middlewares.CheckAuth, controller.GetUserProfile)
}
