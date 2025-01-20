package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bartholomeas/hwheels_api/internal/users/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	api := gin.Default()
	v1 := api.Group("/api/v1")

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API v1 is running"})
	})

	router.InitUsersRouter(v1)

	port := os.Getenv("PORT")
	api.Run(":" + port)
	log.Printf("Server is running on port %s\n", port)
}
