package entities

import (
	"github.com/bartholomeas/hwheels_api/internal/catalog/models"
	commonEntities "github.com/bartholomeas/hwheels_api/internal/common/entities"
)

type CatalogItem struct {
	commonEntities.Base
	Name        string                   `gorm:"not null" json:"name"`
	ModelNumber string                   `json:"model_number"`
	Year        int                      `json:"year"`
	Rarity      models.CatalogItemRarity `json:"rarity"`
	IsChase     bool                     `json:"is_chase" gorm:"default:false"`
	PhotoUrl    string                   `json:"photo_url"`
}
