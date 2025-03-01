package services

import (
	"context"
	"errors"

	authEntities "github.com/bartholomeas/hwheels_api/internal/auth/entities"
	"github.com/bartholomeas/hwheels_api/internal/auth/requests"
	userEntities "github.com/bartholomeas/hwheels_api/internal/user/entities"
	"golang.org/x/crypto/bcrypt"
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

func (s *AuthService) CreateUser(request requests.CreateUserRequest) (*authEntities.User, error) {
	if err := s.checkUserExists(request.Email); err != nil {
		return nil, err
	}

	passwordHash, err := s.hashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	s.cognitoClient.SignUpCognito(context.Background(), request.Username, request.Password, request.Email)

	user := &authEntities.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(passwordHash),
		Role:     authEntities.RoleUser,
	}

	if err := s.createUserWithProfile(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginUser(request requests.LoginRequest) (*string, error) {

	var userFound *authEntities.User

	s.db.Where("email=?", request.Email).Find(&userFound)

	// if userFound.ID == "" {
	// 	return nil, InvalidCredentialsError()
	// }

	// if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(request.Password)); err != nil {
	// 	return nil, InvalidCredentialsError()
	// }

	// generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"id":  userFound.ID,
	// 	"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	// })

	// token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	authResult, err := s.cognitoClient.SignInCognito(context.Background(), request.Email, request.Password)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// return &token, nil
	return authResult.AccessToken, nil
}

func InvalidCredentialsError() error {
	return errors.New("credentials are invalid")
}

//  Helper methods

func (s *AuthService) checkUserExists(email string) error {
	var userFound authEntities.User
	s.db.Where("email = ?", email).First(&userFound)
	if userFound.ID != "" {
		return errors.New("user already exists")
	}
	return nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) createUserWithProfile(user *authEntities.User) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		profile := userEntities.UserProfile{
			UserID: user.ID,
		}

		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		return nil
	})
}
