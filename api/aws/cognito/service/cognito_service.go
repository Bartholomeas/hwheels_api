package cognitoProvider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	appErrors "github.com/bartholomeas/hwheels_api/api/common/app_errors"
)

type CognitoInterface interface {
	SignUpCognito(ctx context.Context, username string, password string, email string) (*string, *appErrors.AppError)
	SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, *appErrors.AppError)
	GetUserByToken(token string) (*cognito.GetUserOutput, *appErrors.AppError)
}

type CognitoService struct {
	cognitoClient *cognito.Client
	clientId      string
	clientSecret  string
}

func NewCognitoService() *CognitoService {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}

	cognitoClient := cognito.NewFromConfig(sdkConfig)

	return &CognitoService{
		cognitoClient: cognitoClient,
		clientId:      os.Getenv("COGNITO_CLIENT_ID"),
		clientSecret:  os.Getenv("COGNITO_CLIENT_SECRET"),
	}
}

func (c CognitoService) SignUpCognito(ctx context.Context, username string, password string, email string) (*string, *appErrors.AppError) {
	secretHash := calculateSecretHash(
		email,
		c.clientId,
		c.clientSecret,
	)

	output, err := c.cognitoClient.SignUp(ctx, &cognito.SignUpInput{
		ClientId:   aws.String(c.clientId),
		Password:   aws.String(password),
		Username:   aws.String(email),
		SecretHash: aws.String(secretHash),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("nickname"),
				Value: aws.String(username),
			},
		},
	})

	if err != nil {
		return nil, parseCognitoError(err)
	}

	return output.UserSub, nil
}

func (c CognitoService) SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, *appErrors.AppError) {
	var authResult *types.AuthenticationResultType

	secretHash := calculateSecretHash(
		username,
		c.clientId,
		c.clientSecret,
	)

	result, err := c.cognitoClient.InitiateAuth(ctx, &cognito.InitiateAuthInput{
		AuthFlow:       "USER_PASSWORD_AUTH",
		ClientId:       aws.String(c.clientId),
		AuthParameters: map[string]string{"USERNAME": username, "PASSWORD": password, "SECRET_HASH": secretHash},
	})

	if err != nil {
		return nil, parseCognitoError(err)
	} else {
		authResult = result.AuthenticationResult
	}

	return authResult, nil
}

func (c CognitoService) GetUserByToken(token string) (*cognito.GetUserOutput, *appErrors.AppError) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(token),
	}

	result, err := c.cognitoClient.GetUser(context.Background(), input)

	if err != nil {
		return nil, parseCognitoError(err)
	}

	return result, nil
}

func parseCognitoError(err error) *appErrors.AppError {
	var notAuthAerr *types.NotAuthorizedException
	var tooManyRequestsErr *types.TooManyRequestsException
	var invalidPasswordErr *types.InvalidPasswordException
	var userNotConfirmedErr *types.UserNotConfirmedException
	var usernameExistsErr *types.UsernameExistsException
	var resetRequired *types.PasswordResetRequiredException
	var tokenExpired *types.ExpiredCodeException

	switch {
	case errors.As(err, &notAuthAerr):
		return appErrors.NewAppError(
			"NotAuthorizedException",
			strings.Split(err.Error(), "NotAuthorizedException: ")[1],
			http.StatusUnauthorized,
		)
	case errors.As(err, &userNotConfirmedErr):
		return appErrors.NewAppError(
			"UserNotConfirmedException",
			strings.Split(err.Error(), "UserNotConfirmedException: ")[1],
			http.StatusBadRequest,
		)
	case errors.As(err, &tooManyRequestsErr):
		return appErrors.NewAppError(
			"TooManyRequestsException",
			strings.Split(err.Error(), "TooManyRequestsException: ")[1],
			http.StatusTooManyRequests,
		)
	case errors.As(err, &invalidPasswordErr):
		return appErrors.NewAppError(
			"InvalidPasswordException",
			strings.Split(err.Error(), "InvalidPasswordException: ")[1],
			http.StatusBadRequest,
		)
	case errors.As(err, &resetRequired):
		return appErrors.NewAppError(
			"PasswordResetRequiredException",
			strings.Split(err.Error(), "PasswordResetRequiredException: ")[1],
			http.StatusBadRequest,
		)
	case errors.As(err, &tokenExpired):
		return appErrors.NewAppError(
			"ExpiredCodeException",
			strings.Split(err.Error(), "ExpiredCodeException: ")[1],
			http.StatusBadRequest,
		)
	case errors.As(err, &usernameExistsErr):
		return appErrors.NewAppError(
			"UsernameExistsException",
			strings.Split(err.Error(), "UsernameExistsException: ")[1],
			http.StatusBadRequest,
		)
	default:
		return appErrors.NewAppError(
			"UnknownError",
			err.Error(),
			http.StatusInternalServerError,
		)
	}

}

func calculateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
