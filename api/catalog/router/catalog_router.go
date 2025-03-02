package router

import (
	controllers "github.com/bartholomeas/hwheels_api/api/catalog/controller"
	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/gin-gonic/gin"
)

func InitCatalogRouter(router *gin.RouterGroup) {
	group := router.Group("/catalog")
	controller := controllers.NewCatalogController(initializers.DB)

	group.GET("/", controller.FindAllItems)
}
