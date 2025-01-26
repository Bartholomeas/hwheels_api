package services

import (
	"errors"
	"fmt"

	"github.com/bartholomeas/hwheels_api/config/initializers"
	authEntities "github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	userEntities "github.com/bartholomeas/hwheels_api/internal/user/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(request requests.AuthInputRequest) (*authEntities.User, error) {
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
		profile := userEntities.Profile{
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
