package basiq

import (
	"time"

	"github.com/basiqio/basiq-sdk-golang/errors"
	"github.com/basiqio/basiq-sdk-golang/utilities"
	"github.com/basiqio/basiq-sdk-golang/v1"
	"github.com/basiqio/basiq-sdk-golang/v2"
)

func NewSessionV1(apiKey string) (*v1.Session, *errors.APIError) {
	session := &v1.Session{
		ApiKey:     apiKey,
		ApiVersion: "1.0",
		Api:        utilities.NewAPI("https://au-api.basiq.io/"),
		Token: &utilities.Token{
			Value:     "",
			Validity:  0,
			Refreshed: time.Now(),
		},
	}

	token, err := utilities.GetToken(apiKey, session.ApiVersion)
	if err != nil {
		return session, err
	}
	session.Token = token
	session.Api.SetHeader("Authorization", "Bearer "+session.Token.Value)

	return session, nil
}

func NewSessionV2(apiKey string) (*v2.Session, *errors.APIError) {
	session := &v2.Session{
		ApiKey:     apiKey,
		ApiVersion: "2.0",
		Api:        utilities.NewAPI("https://au-api.basiq.io/"),
		Token: &utilities.Token{
			Value:     "",
			Validity:  0,
			Refreshed: time.Now(),
		},
	}

	token, err := utilities.GetToken(apiKey, session.ApiVersion)
	if err != nil {
		return session, err
	}
	session.Token = token
	session.Api.SetHeader("Authorization", "Bearer "+session.Token.Value)

	return session, nil
}
