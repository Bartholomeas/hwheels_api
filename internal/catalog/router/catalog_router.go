package router

import (
	"github.com/bartholomeas/hwheels_api/config/initializers"
	controllers "github.com/bartholomeas/hwheels_api/internal/catalog/controller"
	"github.com/gin-gonic/gin"
)

func InitCatalogRouter(router *gin.RouterGroup) {
	group := router.Group("/catalog")
	controller := controllers.NewCatalogController(initializers.DB)

	group.GET("/", controller.FindAllItems)
}
