package controllers

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	"github.com/bartholomeas/hwheels_api/internal/auth/services"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var authRequest requests.AuthInputRequest

	if err := c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(authRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
}
