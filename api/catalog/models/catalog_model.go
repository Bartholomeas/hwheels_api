package models

type CatalogItemRarity string

const (
	CatalogItemRarityCommon CatalogItemRarity = "common"
	CatalogItemRarityRare   CatalogItemRarity = "rare"
	CatalogItemRarityEpic   CatalogItemRarity = "epic"
	CatalogItemRarityLegend CatalogItemRarity = "legend"
)
