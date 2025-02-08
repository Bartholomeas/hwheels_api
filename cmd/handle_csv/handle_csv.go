package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	catalogEntities "github.com/bartholomeas/hwheels_api/internal/catalog/entities"
	"github.com/bartholomeas/hwheels_api/internal/catalog/models"
	"github.com/gosimple/slug"
)

func main() {
	initializers.LoadEnv()
	initializers.ConnectDB()

	startTime := time.Now()

	filePath := flag.String("file", "", "Path to the CSV file")
	flag.Parse()

	if *filePath == "" {
		panic("Please provide a file path using -file flag")
	}

	fmt.Println("Processing file:", *filePath)
	if cwd, err := os.Getwd(); err == nil {
		fmt.Println("Current working directory:", cwd)
	}

	data, err := handleCSV(*filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Data:", data)

	duration := time.Since(startTime)
	fmt.Printf("Execution time: %v\n", duration)

}

func handleCategories(data [][]string) (map[string]struct{}, error) {
	db := initializers.DB

	categoryMap := make(map[string]struct{})
	categoryList := make([]*catalogEntities.CatalogCategory, 0, len(categoryMap))

	for category := range categoryMap {
		categoryList = append(categoryList, &catalogEntities.CatalogCategory{
			Name: category,
		})
	}

	for category := range categoryMap {
		categoryList = append(categoryList, &catalogEntities.CatalogCategory{
			Name:        category,
			Slug:        slug.Make(category),
			Description: "Category description",
		})
	}

	categoriesResult := db.CreateInBatches(categoryList, 100)

	if categoriesResult.Error != nil {
		return nil, fmt.Errorf("failde to batch insert categories: %w", categoriesResult.Error)
	}

	var categoriesWithIds []catalogEntities.CatalogCategory
	if err := db.Find(&categoriesWithIds).Error; err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	return nil, nil
}

func saveToDB(data [][]string) error {
	if len(data) <= 1 {
		return fmt.Errorf("no data to process")
	}

	db := initializers.DB

	categories := make(map[string]struct{})

	for i := 1; i < len(data); i++ {
		row := data[i]
		if len(row) < 5 {
			continue
		}
		categories[row[3]] = struct{}{}
	}

	categoryList := make([]*catalogEntities.CatalogCategory, 0, len(categories))

	for category := range categories {
		categoryList = append(categoryList, &catalogEntities.CatalogCategory{
			Name:        category,
			Slug:        slug.Make(category),
			Description: "Category description",
		})
	}

	if err := db.CreateInBatches(categoryList, 100).Error; err != nil {
		return fmt.Errorf("failed to create categories: %w", err)
	}

	var savedCategories []catalogEntities.CatalogCategory
	if err := db.Find(&savedCategories).Error; err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}

	categoryMap := make(map[string]*catalogEntities.CatalogCategory)
	for i := range savedCategories {
		categoryMap[savedCategories[i].Name] = &savedCategories[i]
	}

	items := make([]*catalogEntities.CatalogItem, 0, len(data)-1)
	for i := 1; i < len(data); i++ {
		row := data[i]
		if len(row) < 5 {
			continue
		}

		category, exists := categoryMap[row[3]]
		if !exists {
			continue
		}

		item := &catalogEntities.CatalogItem{
			Name:        row[2],
			ModelNumber: row[1],
			ReleaseDate: time.Now(),
			RetailPrice: 0,
			MarketValue: 0,
			// Series:      row[3],
			Year:       uint(time.Now().Year()),
			PhotoUrl:   row[4],
			Rarity:     models.CatalogItemRarityCommon,
			IsChase:    false,
			Categories: []*catalogEntities.CatalogCategory{category},
		}
		items = append(items, item)
	}

	if err := db.CreateInBatches(items, 100).Error; err != nil {
		return fmt.Errorf("failed to create items: %w", err)
	}

	return nil
}

func handleCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	data, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	saveToDB(data)

	return data, err
}
