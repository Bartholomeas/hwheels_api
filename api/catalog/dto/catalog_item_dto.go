package dto

import (
	"time"

	"github.com/bartholomeas/hwheels_api/api/catalog/models"
)

type CatalogItemCategoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CatalogItemDTO struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	CreatedAt   time.Time                 `json:"created_at"`
	ModelNumber string                    `json:"model_number"`
	ReleaseDate time.Time                 `json:"release_date"`
	RetailPrice float64                   `json:"retail_price"`
	MarketValue float64                   `json:"market_value"`
	Year        uint                      `json:"year"`
	Rarity      models.CatalogItemRarity  `json:"rarity"`
	IsChase     bool                      `json:"is_chase"`
	PhotoUrl    string                    `json:"photo_url"`
	Categories  []*CatalogItemCategoryDTO `json:"categories"`
}
