package v2

import (
	"encoding/json"
	"fmt"

	"github.com/basiqio/basiq-sdk-golang/errors"
)

type Statement struct {
	ID    string `json:"id"`
	Links struct {
		Self    string `json:"self"`
		Account string `json:"account"`
	} `json:"links"`
}

type StatementList struct {
	Data  []Statement `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type StatementService struct {
	Session *Session
	user    *User
}

func NewStatementService(session *Session, user *User) *StatementService {
	return &StatementService{session, user}
}

func (is *StatementService) GetStatements() (StatementList, *errors.APIError) {
	var data StatementList

	body, _, err := is.Session.Api.Send("GET", "user/"+is.user.Id+"/statements", nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (is *StatementService) GetStatement(statementId string) (Statement, *errors.APIError) {
	var data Statement

	body, _, err := is.Session.Api.Send("GET", "user/"+is.user.Id+"/statements/"+statementId, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}
