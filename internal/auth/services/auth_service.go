package services

import (
	"context"
	"log"

	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	appErrors "github.com/bartholomeas/hwheels_api/internal/common/app_errors"
	"gorm.io/gorm"

	cognito "github.com/bartholomeas/hwheels_api/internal/aws/cognito"
)

type AuthService struct {
	db            *gorm.DB
	cognitoClient cognito.CognitoInterface
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db, cognitoClient: cognito.NewCognitoService()}
}

func (s *AuthService) CreateUser(request requests.CreateUserRequest) (any, error) {

	log.Println("REQUEST:: ", request)
	createdUser, err := s.cognitoClient.SignUpCognito(context.Background(), request.Username, request.Password, request.Email)

	if err != nil {
		return nil, err
	}
	log.Println("ress", createdUser)

	return createdUser, nil
}

func (s *AuthService) LoginUser(request requests.LoginRequest) (*string, error) {

	authResult, err := s.cognitoClient.SignInCognito(context.Background(), request.Email, request.Password)
	log.Printf("ERROR::::: %v", err)
	if err != nil {
		return nil, appErrors.NewAppError(err.Code, err.Message, 400)
	}

	return authResult.AccessToken, nil
}
