package controllers

import (
	"net/http"
	"strconv"

	"github.com/bartholomeas/hwheels_api/api/catalog/repositories"
	"github.com/bartholomeas/hwheels_api/api/catalog/services"
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
	page := 1
	limit := 24

	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	items, err := c.catalogService.FindAll(limit, page)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}
