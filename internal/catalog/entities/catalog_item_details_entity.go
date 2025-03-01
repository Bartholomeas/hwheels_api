package entities

import commonEntities "github.com/bartholomeas/hwheels_api/internal/common/entities"

type CatalogItemDetails struct {
	commonEntities.Base
	Description   string       `json:"description"`
	CatalogItemID string       `json:"catalog_item_id" gorm:"type:uuid;not null;unique"`
	CatalogItem   *CatalogItem `json:"catalog_item" gorm:"type:uuid;not null;unique"`
}

func (CatalogItemDetails) TableName() string {
	return "catalog_item_details"
}
