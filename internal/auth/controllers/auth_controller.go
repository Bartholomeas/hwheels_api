package controllers

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput requests.CreateUserRequest

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound entities.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{
		Username: authInput.Username,
		Password: string(passwordHash),
		Email:    authInput.Email,
		Role:     "user",
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
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
