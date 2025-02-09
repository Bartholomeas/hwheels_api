package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	authRouter "github.com/bartholomeas/hwheels_api/internal/auth/router"
	catalogRouter "github.com/bartholomeas/hwheels_api/internal/catalog/router"
	userRouter "github.com/bartholomeas/hwheels_api/internal/user/router"
	"github.com/gin-gonic/gin"
)

func main() {

	initializers.LoadEnv()
	initializers.ConnectDB()

	api := gin.Default()
	v1 := api.Group("/api/v1")

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API v1 is running"})
	})

	authRouter.InitAuthRouter(v1)
	userRouter.InitUserRouter(v1)
	catalogRouter.InitCatalogRouter(v1)

	port := os.Getenv("PORT")
	api.Run(":" + port)
	log.Printf("Server is running on port %s\n", port)
}
