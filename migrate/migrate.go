package main

import (
	"github.com/bartholomeas/hwheels_api/config/initializers"
	models "github.com/bartholomeas/hwheels_api/internal/auth/model"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
}
