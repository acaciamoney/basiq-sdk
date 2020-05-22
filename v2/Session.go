package v2

import (
	"github.com/basiqio/basiq-sdk-golang/errors"
	"github.com/basiqio/basiq-sdk-golang/utilities"
)

type Session struct {
	ApiKey     string
	ApiVersion string
	Api        *utilities.API
	Token      *utilities.Token
}

func (s *Session) RefreshToken() *errors.APIError {
	token, err := utilities.GetToken(s.ApiKey, s.ApiVersion)
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
