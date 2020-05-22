package utilities

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/basiqio/basiq-sdk-golang/errors"
)

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

	return &Token{
		Value:     data.AccessToken,
		Validity:  time.Duration(data.ExpiresIn) * time.Second,
		Refreshed: time.Now(),
	}, nil
}
