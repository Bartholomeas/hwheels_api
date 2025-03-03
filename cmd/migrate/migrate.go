package main

import (
	"log"

	authEntities "github.com/bartholomeas/hwheels_api/api/auth/entities"
	catalogEntities "github.com/bartholomeas/hwheels_api/api/catalog/entities"
	userEntities "github.com/bartholomeas/hwheels_api/api/user/entities"
	"github.com/bartholomeas/hwheels_api/config/initializers"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	log.Println("Starting migrations...")

	if err := initializers.DB.Migrator().DropTable(
		"catalog_item_categories",
		// &catalogEntities.CatalogItemDetails{},
		// &catalogEntities.CatalogItem{},
		// &catalogEntities.CatalogCategory{},
		&userEntities.UserProfile{},
		&authEntities.User{},
	); err != nil {
		log.Fatal("Error dropping tables:", err)
	}

	if err := initializers.DB.AutoMigrate(
		&authEntities.User{},
		&userEntities.UserProfile{},
		&catalogEntities.CatalogCategory{},
		&catalogEntities.CatalogItem{},
		&catalogEntities.CatalogItemDetails{},
	); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
