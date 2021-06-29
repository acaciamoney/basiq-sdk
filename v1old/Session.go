package v1old

import (
	"encoding/json"
	"fmt"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/acaciamoney/basiq-sdk/utilities"
)

type Session struct {
	ApiKey     string
	ApiVersion string
	Api        *utilities.API
	Token      *utilities.Token
}

func (s *Session) RefreshToken() *errors.APIError {
	token, err := utilities.GetToken(s.ApiKey, "1.0")
	if err != nil {
		return err
	}
	s.Token = token
	s.Api.SetHeader("Authorization", "Bearer "+token.Value)
	return nil
}

func (s *Session) CreateUser(createData *UserData) (User, *errors.APIError) {
	return NewUserService(s).CreateUser(createData)
}

func (s *Session) ForUser(userId string) User {
	return NewUserService(s).ForUser(userId)
}

func (s *Session) GetInstitutions() (InstitutionsList, *errors.APIError) {
	return NewInstitutionService(s).GetInstitutions()
}

func (s *Session) GetInstitution(id string) (Institution, *errors.APIError) {
	return NewInstitutionService(s).GetInstitution(id)
}

func (s *Session) GetJob(id string) (Job, *errors.APIError) {
	var data Job

	body, _, err := s.Api.Send("GET", "jobs/"+id, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}
