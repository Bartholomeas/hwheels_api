package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	authEntities "github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	userEntities "github.com/bartholomeas/hwheels_api/internal/user/entities"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(request requests.CreateUserRequest) (*authEntities.User, error) {
	var userFound authEntities.User

	initializers.DB.Where("email=?", request.Email).First(&userFound)
	if userFound.ID != "" {
		return nil, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := authEntities.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(passwordHash),
		Role:     authEntities.RoleUser,
	}

	err = initializers.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		fmt.Println("JUZERRR:: ", user.ID)
		profile := userEntities.UserProfile{
			UserID: user.ID,
		}

		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(request requests.LoginRequest) (*string, error) {

	var userFound *authEntities.User

	initializers.DB.Where("email=?", request.Email).Find(&userFound)

	if userFound.ID == "" {
		return nil, InvalidCredentialsError()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(request.Password)); err != nil {
		return nil, InvalidCredentialsError()
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &token, nil
}

func InvalidCredentialsError() error {
	return errors.New("credentials are invalid")
}
