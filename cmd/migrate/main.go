package main

import (
	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/bartholomeas/hwheels_api/internal/users/model"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&model.User{})
} 