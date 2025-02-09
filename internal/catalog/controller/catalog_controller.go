package controllers

import (
	"net/http"

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

func (c *CatalogController) FindAllItems(ctx *gin.Context) {
	items, err := c.catalogService.FindAllItems()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}
