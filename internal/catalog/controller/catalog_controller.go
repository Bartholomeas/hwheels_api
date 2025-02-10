package controllers

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/internal/catalog/repositories"
	"github.com/bartholomeas/hwheels_api/internal/catalog/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CatalogController struct {
	catalogService *services.CatalogService
}

func NewCatalogController(db *gorm.DB) *CatalogController {
	repo := repositories.NewCatalogRepository(db)
	service := services.NewCatalogService(repo)

	return &CatalogController{
		catalogService: service,
	}
}

func (c *CatalogController) FindAllItems(ctx *gin.Context) {
	items, err := c.catalogService.FindAll(10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}
