package repositories

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"gorm.io/gorm"
)

type CatalogRepository interface {
	FindAll() ([]entities.CatalogItem, error)
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{db}
}

func (r *catalogRepository) FindAll() ([]entities.CatalogItem, error) {

	var items []entities.CatalogItem

	err := r.db.
		Model(&entities.CatalogItem{}).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug")
		}).
		Select("catalog_items.*").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
