package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type CachedToken struct {
	Token  string
	Expiry int
}

type Token struct {
	Value     string
	Validity  time.Duration
	Refreshed time.Time
}

type AuthorizationResponse struct {
	AccessToken string        `json:"access_token"`
	Type        string        `json:"type"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

func GetToken(apiKey, apiVersion string) (*Token, *errors.APIError) {

	token := GetCachedToken()
	log.Print(token)
	if token.Expiry > int(time.Now().Unix()) {

		return &Token{
			Value:     token.Token,
			Validity:  time.Duration(token.Expiry) * time.Second,
			Refreshed: time.Now(),
		}, nil
	}
	log.Print("TOKEN HAS EXPIRED!!")
	body, _, err := NewAPI("https://au-api.basiq.io/").SetHeader("Authorization", "Basic "+apiKey).
		SetHeader("basiq-version", apiVersion).
		SetHeader("content-type", "application/json").
		Send("POST", "token", nil)
	if err != nil {
		return nil, err
	}

	var data AuthorizationResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return nil, &errors.APIError{Message: err.Error()}
	}
	expiry := time.Duration(data.ExpiresIn) * time.Second

	SetCachedToken(CachedToken{
		Token:  data.AccessToken,
		Expiry: int(time.Now().Unix()) + int(expiry.Seconds()),
	})

	return &Token{
		Value:     data.AccessToken,
		Validity:  time.Duration(data.ExpiresIn) * time.Second,
		Refreshed: time.Now(),
	}, nil
}

func GetCachedToken() CachedToken {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	payload, err := secretsmanager.New(sess).GetSecretValue(
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String("BasiqToken"),
		},
	)
	if err != nil {
		log.Print("Unable to fetch cached basiq token")
		log.Fatal(err)
	}
	var token CachedToken
	err = json.Unmarshal([]byte(*payload.SecretString), &token)
	if err != nil {
		log.Print("Unable to parse cached basiq token")
		log.Fatal(err)
	}
	return token
}

func SetCachedToken(t CachedToken) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	str, err := json.Marshal(t)
	if err != nil {
		log.Print("Unable to serialise cached basiq token")
		log.Fatal(err)
	}
	_, err = secretsmanager.New(sess).PutSecretValue(
		&secretsmanager.PutSecretValueInput{
			SecretId:     aws.String("BasiqToken"),
			SecretString: aws.String(string(str)),
		},
	)
	if err != nil {
		log.Print("Unable Save cached basiq token...")
		log.Fatal(err)
	}
}
