package main

import (
	"log"
	"net/http"
	"os"

	authRouter "github.com/bartholomeas/hwheels_api/api/auth/router"
	catalogRouter "github.com/bartholomeas/hwheels_api/api/catalog/router"
	userRouter "github.com/bartholomeas/hwheels_api/api/user/router"
	"github.com/bartholomeas/hwheels_api/config/initializers"
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
