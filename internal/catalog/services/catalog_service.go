package services

import (
	catalogDto "github.com/bartholomeas/hwheels_api/internal/catalog/dto"
	"github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"github.com/bartholomeas/hwheels_api/internal/catalog/repositories"
	"github.com/bartholomeas/hwheels_api/internal/common/pagination"
	paginationDto "github.com/bartholomeas/hwheels_api/internal/common/pagination/dto"
)

type CatalogService struct {
	repo repositories.CatalogRepository
}

func NewCatalogService(repo repositories.CatalogRepository) *CatalogService {
	return &CatalogService{repo}
}

// func (s *CatalogService) FindAll(limit int, offset int) ([]dto.CatalogItemDTO, error) {
func (s *CatalogService) FindAll(limit int, page int) (*paginationDto.PaginatedResponse[catalogDto.CatalogItemDTO], error) {

	paginationParams := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	params := repositories.FindAllParams{
		Limit:  paginationParams.GetLimit(),
		Offset: paginationParams.GetOffset(),
	}

	items, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, err
	}

	paginationParams.ItemsCount = total
	paginationParams.PagesCount = int(total) / paginationParams.GetLimit()
	if int(total)%paginationParams.GetLimit() > 0 {
		paginationParams.PagesCount++
	}

	meta := paginationDto.Metadata{
		Page:        paginationParams.Page,
		Limit:       paginationParams.Limit,
		TotalCount:  int(total),
		TotalPages:  paginationParams.PagesCount,
		HasNextPage: paginationParams.GetPage() < paginationParams.PagesCount,
		HasPrevPage: paginationParams.GetPage() > 1,
	}

	itemDTOs := make([]catalogDto.CatalogItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = mapToCatalogItemDTO(item)
	}

	response := &paginationDto.PaginatedResponse[catalogDto.CatalogItemDTO]{
		Meta: meta,
		Data: itemDTOs,
	}

	return response, nil

}

// Private stuff

func mapToCatalogItemDTO(item entities.CatalogItem) catalogDto.CatalogItemDTO {
	categories := make([]*catalogDto.CatalogItemCategoryDTO, len(item.Categories))

	for j, category := range item.Categories {
		categories[j] = &catalogDto.CatalogItemCategoryDTO{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		}
	}
	return catalogDto.CatalogItemDTO{
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
