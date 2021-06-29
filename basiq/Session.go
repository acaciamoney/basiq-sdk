package basiq

import (
	"time"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/acaciamoney/basiq-sdk/utilities"
	v2 "github.com/acaciamoney/basiq-sdk/v2"
)

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
