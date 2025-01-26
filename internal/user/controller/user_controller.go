package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	user, _ := c.Get("currentUser")

	c.JSON(http.StatusOK, gin.H{"user": user})
}
