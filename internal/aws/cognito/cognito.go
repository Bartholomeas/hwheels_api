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
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoInterface interface {
	SignUpCognito(ctx context.Context, username string, password string, email string) (bool, error)
	SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, error)
}

type CognitoService struct {
	CognitoClient *cognitoidentityprovider.Client
	clientId      string
	clientSecret  string
}

func NewCognitoService() *CognitoService {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}

	cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)

	return &CognitoService{
		CognitoClient: cognitoClient,
		clientId:      os.Getenv("COGNITO_CLIENT_ID"),
		clientSecret:  os.Getenv("COGNITO_CLIENT_SECRET"),
	}
}

func calculateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (actor CognitoService) SignUpCognito(ctx context.Context, username string, password string, email string) (bool, error) {
	confirmed := false

	secretHash := calculateSecretHash(
		email,
		actor.clientId,
		actor.clientSecret,
	)

	output, err := actor.CognitoClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(actor.clientId),
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

func (actor CognitoService) SignInCognito(ctx context.Context, username string, password string) (*types.AuthenticationResultType, error) {
	var authResult *types.AuthenticationResultType

	secretHash := calculateSecretHash(
		username,
		actor.clientId,
		actor.clientSecret,
	)

	result, err := actor.CognitoClient.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       "USER_PASSWORD_AUTH",
		ClientId:       aws.String(actor.clientId),
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
