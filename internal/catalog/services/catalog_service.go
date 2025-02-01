package services

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"gorm.io/gorm"
)

type CatalogService struct {
	db *gorm.DB
}

func NewCatalogService(db *gorm.DB) *CatalogService {
	return &CatalogService{db: db}
}

func (s *CatalogService) FindAllItems() ([]entities.CatalogItem, error) {
	// var items []entities.CatalogItem
	// if err := s.db.Find(&items).Error; err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
