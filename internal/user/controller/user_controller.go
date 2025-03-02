package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	"github.com/bartholomeas/hwheels_api/internal/auth/entities"
	cognito "github.com/bartholomeas/hwheels_api/internal/aws/cognito"
	appErrors "github.com/bartholomeas/hwheels_api/internal/common/app_errors"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	GetUserProfile(c *gin.Context)
	GetUserByToken(c *gin.Context) (*UserResponse, error)
}

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type UserController struct {
	cognitoClient cognito.CognitoInterface
}

func NewUserController() *UserController {
	return &UserController{
		cognitoClient: cognito.NewCognitoService(),
	}
}

func (actor *UserController) GetUserProfile(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user, ok := currentUser.(entities.User)

	if !ok {
		c.JSON(http.StatusBadRequest, *appErrors.NewAppError("Could not get user from context", "Could not get user from context", http.StatusBadRequest))
		return
	}

	var userWithProfile entities.User
	initializers.DB.
		Preload("Profile").
		Where("id = ?", user.ID).
		First(&userWithProfile)

	c.JSON(http.StatusOK, userWithProfile)
}

func (actor *UserController) GetUserByToken(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	if token == "" {
		var tokenError = appErrors.NewAppError("Token not found", "Token not found", http.StatusUnauthorized)

		c.JSON(tokenError.StatusCode, tokenError)
		return
	}

	cognitoUser, err := actor.cognitoClient.GetUserByToken(token)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	user := &UserResponse{}
	for _, attribute := range cognitoUser.UserAttributes {
		switch *attribute.Name {
		case "sub":
			user.ID = *attribute.Value
		case "nickname":
			user.Username = *attribute.Value
		case "email":
			user.Email = *attribute.Value
		case "email_verified":
			emailVerified, err := strconv.ParseBool(*attribute.Value)
			if err == nil {
				user.EmailVerified = emailVerified
			}
		}
	}

	c.JSON(http.StatusOK, user)

}
