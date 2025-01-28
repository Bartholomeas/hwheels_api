package controller

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user, ok := currentUser.(entities.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user from context"})
		return
	}

	var userWithProfile entities.User
	initializers.DB.
		Preload("Profile").
		Where("id = ?", user.ID).
		First(&userWithProfile)

	c.JSON(http.StatusOK, userWithProfile)
}
