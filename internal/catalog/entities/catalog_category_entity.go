package entities

import "github.com/bartholomeas/hwheels_api/internal/common/entities"

type CatalogCategory struct {
	entities.Base
	Name        string         `json:"name" gorm:"not null"`
	Slug        string         `json:"slug" gorm:"not null"`
	Description string         `json:"description"`
	Items       []*CatalogItem `json:"items" gorm:"many2many:catalog_item_categories;"`
}

func (CatalogCategory) TableName() string {
	return "catalog_categories"
}
