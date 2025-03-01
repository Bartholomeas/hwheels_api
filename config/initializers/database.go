package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	connString := os.Getenv("DB_CONNECTION_STRING")
	log.Printf("Connecting to database with: %s", connString)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(connString), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")
}
