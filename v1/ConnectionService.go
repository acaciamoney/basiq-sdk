package v1

import (
	"encoding/json"
	"fmt"

	"github.com/basiqio/basiq-sdk-golang/errors"
)

type ConnectionService struct {
	Session Session
	user    User
}

type Connection struct {
	Id          string       `json:"id"`
	Status      string       `json:"status"`
	LastUsed    string       `json:"lastUsed"`
	Institution Institution  `json:"institution"`
	Accounts    AccountsList `json:"accounts"`
	Service     *ConnectionService
}

type InstitutionData struct {
	Id string `json:"id"`
}

type ConnectionData struct {
	Institution  *InstitutionData `json:"institution"`
	LoginId      string           `json:"loginId"`
	Password     string           `json:"password"`
	SecurityCode string           `json:"securityCode,omitempty"`
}

type ConnectionFilter struct {
	Id            string `json:"id,omitempty"`
	Status        string `json:"status,omitempty"`
	InstitutionId string `json:"institution.id,omitempty"`
}

func NewConnectionService(session *Session, user *User) *ConnectionService {
	return &ConnectionService{
		Session: *session,
		user:    *user,
	}
}

func (cs *ConnectionService) GetConnection(connectionId string) (Connection, *errors.APIError) {
	var data Connection

	data.Service = cs

	body, _, err := cs.Session.Api.Send("GET", "users/"+cs.user.Id+"/connections/"+connectionId, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (cs *ConnectionService) ForConnection(connectionId string) Connection {
	var data Connection

	data.Service = cs
	data.Id = connectionId

	return data
}

func (cs *ConnectionService) NewConnection(connectionData *ConnectionData) (Job, *errors.APIError) {
	var data Job
	data.Service = cs

	jsonBody, errorr := json.Marshal(connectionData)
	if errorr != nil {
		return data, &errors.APIError{Message: errorr.Error()}
	}

	body, _, err := cs.Session.Api.Send("POST", "users/"+cs.user.Id+"/connections", jsonBody)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (cs *ConnectionService) RefreshConnection(connectionId string) (Job, *errors.APIError) {
	var data Job
	data.Service = cs

	body, _, err := cs.Session.Api.Send("POST", "users/"+cs.user.Id+"/connections/"+connectionId+"/refresh", nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (cs *ConnectionService) UpdateConnection(connectionId, password string) (Job, *errors.APIError) {
	var data Job
	data.Service = cs

	jsonBody := []byte(`{"password":"` + password + `"}`)

	body, _, err := cs.Session.Api.Send("POST", "users/"+cs.user.Id+"/connections/"+connectionId, jsonBody)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (cs *ConnectionService) DeleteConnection(connectionId string) *errors.APIError {
	var data Job
	data.Service = cs

	_, _, err := cs.Session.Api.Send("DELETE", "users/"+cs.user.Id+"/connections/"+connectionId, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ConnectionService) GetJob(jobId string) (Job, *errors.APIError) {
	var data Job
	data.Service = cs

	body, _, err := cs.Session.Api.Send("GET", "jobs/"+jobId, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (c *Connection) Refresh() (Job, *errors.APIError) {
	return c.Service.RefreshConnection(c.Id)
}

func (c *Connection) Update(password string) (Job, *errors.APIError) {
	return c.Service.UpdateConnection(c.Id, password)
}

func (c *Connection) Delete() *errors.APIError {
	return c.Service.DeleteConnection(c.Id)
}
