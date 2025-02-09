package services

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/dto"
	"github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"gorm.io/gorm"
)

type CatalogService struct {
	db *gorm.DB
}

func NewCatalogService(db *gorm.DB) *CatalogService {
	return &CatalogService{db: db}
}

func (s *CatalogService) FindAllItems() ([]dto.CatalogItemDTO, error) {
	var items []entities.CatalogItem

	if err := s.db.
		Model(&entities.CatalogItem{}).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug")
		}).
		Select("catalog_items.*").
		Find(&items).Error; err != nil {
		return nil, err
	}

	itemDTOs := make([]dto.CatalogItemDTO, len(items))

	for i, item := range items {
		categories := make([]*dto.CatalogItemCategoryDTO, len(item.Categories))
		for j, category := range item.Categories {
			categories[j] = &dto.CatalogItemCategoryDTO{
				ID:   category.ID,
				Name: category.Name,
				Slug: category.Slug,
			}
		}

		itemDTOs[i] = dto.CatalogItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			CreatedAt:   item.CreatedAt,
			ModelNumber: item.ModelNumber,
			ReleaseDate: item.ReleaseDate,
			RetailPrice: item.RetailPrice,
			MarketValue: item.MarketValue,
			Year:        item.Year,
			Rarity:      item.Rarity,
			IsChase:     item.IsChase,
			PhotoUrl:    item.PhotoUrl,
			Categories:  categories,
		}
	}

	return itemDTOs, nil

}
