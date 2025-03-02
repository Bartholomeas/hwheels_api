package repositories

import (
	"github.com/bartholomeas/hwheels_api/api/catalog/entities"
	"gorm.io/gorm"
)

type CatalogRepository interface {
	FindAll(params FindAllParams) ([]entities.CatalogItem, int64, error)
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{db}
}

type FindAllParams struct {
	Limit  int
	Offset int
}

func (p *FindAllParams) setDefaults() {
	if p.Limit == 0 {
		p.Limit = 24
	}
	if p.Offset == 0 {
		p.Offset = 1
	}
}

func (r *catalogRepository) FindAll(params FindAllParams) ([]entities.CatalogItem, int64, error) {
	params.setDefaults()

	var items []entities.CatalogItem
	var totalRows int64

	if err := r.db.
		Model(&entities.CatalogItem{}).
		Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Model(&entities.CatalogItem{}).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug")
		}).
		Select("catalog_items.*").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, int64(len(items)), nil
}
