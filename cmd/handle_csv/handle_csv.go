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

func saveToDB(data [][]string) error {
	if len(data) <= 1 {
		return fmt.Errorf("no data to process")
	}

	items := make([]*catalogEntities.CatalogItem, 0, len(data)-1)

	for i := 1; i < len(data); i++ {
		row := data[i]
		if len(row) < 5 {
			continue
		}
		item := &catalogEntities.CatalogItem{
			Name:        row[2],
			ModelNumber: row[1],
			ReleaseDate: time.Now(),
			RetailPrice: 0,
			MarketValue: 0,
			Series:      row[3],
			Year:        uint(time.Now().Year()),
			PhotoUrl:    row[4],
			Rarity:      models.CatalogItemRarityCommon,
			IsChase:     false,
		}
		items = append(items, item)
	}

	db := initializers.DB
	result := db.CreateInBatches(items, 100)

	if result.Error != nil {
		return fmt.Errorf("failde to batch insert items: %w", result.Error)
	}

	// for _, row := range data {
	// 	fmt.Println(row)
	// 	// db.Create(&catalogEntities.CatalogItem{
	// 	// 	Name: row[0],
	// 	// })
	// }

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
	// db := initializers.DB
	saveToDB(data)
	// for _, row := range data {
	// 	for _, col := range row {
	// 		// db.Create(&catalogEntities.CatalogItem{
	// 		// 	Name:        col,
	// 		// 	ModelNumber: col,
	// 		// 	ReleaseDate: time.Now(),
	// 		// 	RetailPrice: 0,
	// 		// 	MarketValue: 0,
	// 		// 	Series:      col,
	// 		// 	Year:        0,
	// 		// 	Rarity:      models.CatalogItemRarityCommon,
	// 		// })

	// 		// fmt.Printf("%s,", col)
	// 	}
	// 	fmt.Println()
	// }

	return data, err
}
