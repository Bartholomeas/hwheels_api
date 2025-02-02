package entities

import commonEntities "github.com/bartholomeas/hwheels_api/internal/common/entities"

type CatalogCategory struct {
	commonEntities.Base
	Name  string         `gorm:"not null" json:"name"`
	Items []*CatalogItem `json:"items" gorm:"many2many:item_categories"`
	// Items []CatalogItem `json:"items" gorm:"foreignKey:CategoryID"`
}
