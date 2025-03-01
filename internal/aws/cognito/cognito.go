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
	CognitoSignUp(ctx context.Context, clientId string, username string, password string, email string) (bool, error)
}

type CognitoService struct {
	CognitoClient *cognitoidentityprovider.Client
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
	}
}

func calculateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (actor CognitoService) CognitoSignUp(ctx context.Context, clientId string, username string, password string, email string) (bool, error) {
	confirmed := false

	secretHash := calculateSecretHash(
		email,
		clientId,
		os.Getenv("COGNITO_CLIENT_SECRET"),
	)

	output, err := actor.CognitoClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientId),
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
