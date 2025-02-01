package controllers

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CatalogController struct {
	catalogService *services.CatalogService
}

func NewCatalogController(db *gorm.DB) *CatalogController {
	return &CatalogController{
		catalogService: services.NewCatalogService(db),
	}
}

func (c *CatalogController) FindAllItems(ctx *gin.Context) {}
