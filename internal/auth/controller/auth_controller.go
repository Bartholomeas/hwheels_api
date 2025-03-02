package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	"github.com/bartholomeas/hwheels_api/internal/auth/services"
	appErrors "github.com/bartholomeas/hwheels_api/internal/common/app_errors"
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

		log.Printf("ERROR:: %v", err.Error())
		appError := *appErrors.NewAppError("InvalidRequest", strings.SplitAfter(err.Error(), "Error:")[1], http.StatusBadRequest)
		ctx.JSON(appError.StatusCode, appError)
		return
	}

	user, err := c.authService.CreateUser(authRequest)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

func (c *AuthController) LoginUser(ctx *gin.Context) {
	var loginRequest requests.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	token, err := c.authService.LoginUser(loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in", "accessToken": token})
}
