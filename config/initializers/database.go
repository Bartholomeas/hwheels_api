package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	connString := os.Getenv("DB_CONNECTION_STRING")
	var err error
	DB, err = gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
}
