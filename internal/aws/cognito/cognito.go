package cognitoProvider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoInterface interface {
	SignUpCognito(ctx context.Context, username string, password string, email string) (bool, error)
	SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, error)
	GetUserByToken(token string) (*cognito.GetUserOutput, error)
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

func calculateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c CognitoService) SignUpCognito(ctx context.Context, username string, password string, email string) (bool, error) {
	confirmed := false

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
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			log.Println(*invalidPassword.Message)
		} else {
			log.Printf("Couldn't sign up user %v. Here's why: %v\n", username, err)
		}

	} else {
		confirmed = output.UserConfirmed
	}

	return confirmed, nil
}

func (c CognitoService) SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, error) {
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
		var resetRequired *types.PasswordResetRequiredException
		if errors.As(err, &resetRequired) {
			log.Println(*resetRequired.Message)
		} else {
			log.Printf("Couldn't sign in user %v. Here's why: %v\n", username, err)
		}
	} else {
		authResult = result.AuthenticationResult
	}

	return authResult, err
}

func (c CognitoService) GetUserByToken(token string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(token),
	}

	result, err := c.cognitoClient.GetUser(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
