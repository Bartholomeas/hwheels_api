package entities

import (
	"time"

	"github.com/bartholomeas/hwheels_api/api/catalog/models"
	commonEntities "github.com/bartholomeas/hwheels_api/api/common/entities"
)

type CatalogItem struct {
	commonEntities.Base
	Name        string                   `gorm:"not null" json:"name"`
	ModelNumber string                   `json:"model_number"`
	ReleaseDate time.Time                `json:"release_date"`
	RetailPrice float64                  `json:"retail_price"`
	MarketValue float64                  `json:"market_value"`
	Year        uint                     `json:"year"`
	Rarity      models.CatalogItemRarity `json:"rarity"`
	IsChase     bool                     `json:"is_chase" gorm:"default:false"`
	PhotoUrl    string                   `json:"photo_url"`
	Categories  []*CatalogCategory       `json:"categories,omitempty" gorm:"many2many:catalog_item_categories;constraint:OnDelete:CASCADE"`
	Details     *CatalogItemDetails      `json:"details,omitempty" gorm:"foreignKey:CatalogItemID;constraint:OnDelete:CASCADE"`
}

func (CatalogItem) TableName() string {
	return "catalog_items"
}
