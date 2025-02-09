package services

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/dto"
	"github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"github.com/bartholomeas/hwheels_api/internal/catalog/repositories"
)

type CatalogService struct {
	repo repositories.CatalogRepository
}

func NewCatalogService(repo repositories.CatalogRepository) *CatalogService {
	return &CatalogService{repo}
}

func (s *CatalogService) FindAll() ([]dto.CatalogItemDTO, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	itemDTOs := make([]dto.CatalogItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = mapToCatalogItemDTO(item)
	}
	return itemDTOs, nil

}

// Private stuff

func mapToCatalogItemDTO(item entities.CatalogItem) dto.CatalogItemDTO {
	categories := make([]*dto.CatalogItemCategoryDTO, len(item.Categories))

	for j, category := range item.Categories {
		categories[j] = &dto.CatalogItemCategoryDTO{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		}
	}
	return dto.CatalogItemDTO{
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
