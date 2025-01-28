package controllers

import (
	"net/http"

	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	"github.com/bartholomeas/hwheels_api/internal/auth/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(db),
	}
}

func (c *AuthController) CreateUser(ctx *gin.Context) {
	var authRequest requests.CreateUserRequest

	if err := ctx.ShouldBindJSON(&authRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.CreateUser(authRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (c *AuthController) LoginUser(ctx *gin.Context) {
	var loginRequest requests.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := c.authService.LoginUser(loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in", "accessToken": token})
}
