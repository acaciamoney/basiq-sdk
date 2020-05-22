package v2

import (
	"encoding/json"
	"fmt"

	"github.com/basiqio/basiq-sdk-golang/errors"
	"github.com/basiqio/basiq-sdk-golang/utilities"
)

type AccountsList struct {
	Count int       `json:"count"`
	Data  []Account `json:"data"`
}

type Account struct {
	Id             string `json:"id"`
	AccountNo      string `json:"accountNo"`
	Name           string `json:"name"`
	Currency       string `json:"currency"`
	Balance        string `json:"balance"`
	AvailableFunds string `json:"availableFunds"`
	LastUpdated    string `json:"lastUpdated"`
	Class          struct {
		Type    string `json:"type"`
		Product string `json:"product"`
		Meta    struct {
			AccountNumber       string `json:"accountNumber"`
			AvailableRedraw     string `json:"availableRedraw"`
			EndDate             string `json:"endDate"`
			Fee                 string `json:"fee"`
			InstalmentAmount    string `json:"instalmentAmount"`
			InterestRate        string `json:"interestRate"`
			InterestType        string `json:"interestType"`
			NextInstalmentDate  string `json:"nextInstalmentDate"`
			OffsetAccountNumber string `json:"offsetAccountNumber"`
			RepaymentFrequency  string `json:"repaymentFrequency"`
			RepaymentType       string `json:"repaymentType"`
		} `json:"meta"`
	} `json:"class"`
	Status      string `json:"status"`
	Institution string `json:"institution"`
	Connection  string `json:"connection"`
}

type UserService struct {
	Session Session
}

type UserData struct {
	Mobile string `json:"mobile,omitempty"`
	Email  string `json:"email,omitempty"`
}

func NewUserService(session *Session) *UserService {
	return &UserService{
		Session: *session,
	}
}

func (us *UserService) CreateUser(createData *UserData) (User, *errors.APIError) {
	var data User

	jsonBody, errorr := json.Marshal(createData)
	if errorr != nil {
		return data, &errors.APIError{Message: errorr.Error()}
	}

	body, _, err := us.Session.Api.Send("POST", "users", jsonBody)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	data.Service = us

	return data, nil
}

func (us *UserService) GetUser(userId string) (User, *errors.APIError) {
	var data User

	body, _, err := us.Session.Api.Send("GET", "users/"+userId, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	data.Service = us

	return data, nil
}

func (us *UserService) UpdateUser(userId string, updateData *UserData) (User, *errors.APIError) {
	var data User

	jsonBody, errorr := json.Marshal(updateData)
	if errorr != nil {
		return data, &errors.APIError{Message: errorr.Error()}
	}

	body, _, err := us.Session.Api.Send("POST", "users/"+userId, jsonBody)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	data.Service = us

	return data, nil
}

func (us *UserService) DeleteUser(userId string) *errors.APIError {
	_, _, err := us.Session.Api.Send("DELETE", "users/"+userId, nil)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) RefreshAllConnections(userId string) *errors.APIError {
	_, _, err := us.Session.Api.Send("POST", "users/"+userId+"/connections/refresh", nil)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) ListAllConnections(userId string, filter *utilities.FilterBuilder) (ConnectionList, *errors.APIError) {
	var data ConnectionList

	url := "users/" + userId + "/connections"
	if filter != nil {
		url = url + "?" + filter.GetFilter()
	}

	body, _, err := us.Session.Api.Send("GET", url, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (us *UserService) GetAccounts(userId string, filter *utilities.FilterBuilder) (AccountsList, *errors.APIError) {
	var data AccountsList

	url := "users/" + userId + "/accounts"

	if filter != nil {
		url = url + "?" + filter.GetFilter()
	}

	body, _, err := us.Session.Api.Send("GET", url, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (us *UserService) GetAccount(userId string, accountId string) (Account, *errors.APIError) {
	var data Account

	body, _, err := us.Session.Api.Send("GET", "users/"+userId+"/accounts/"+accountId, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	return data, nil
}

func (us *UserService) GetTransactions(userId string, filter *utilities.FilterBuilder) (TransactionsList, *errors.APIError) {
	return NewTransactionService(&us.Session).GetTransactions(userId, filter)
}

func (us *UserService) GetTransaction(userId string, transactionId string) (Transaction, *errors.APIError) {
	return NewTransactionService(&us.Session).GetTransaction(userId, transactionId)
}

func (us *UserService) ForUser(userId string) User {
	return User{
		Id:      userId,
		Service: us,
	}
}

type ConnectionList struct {
	Count int          `json:"count"`
	Data  []Connection `json:"data"`
}

type User struct {
	Id          string         `json:"id"`
	Email       string         `json:"email"`
	Mobile      string         `json:"mobile"`
	Connections ConnectionList `json:"connections"`
	Service     *UserService
}

func (u *User) CreateConnection(connectionData *ConnectionData) (Job, *errors.APIError) {
	return NewConnectionService(&u.Service.Session, u).NewConnection(connectionData)
}

func (u *User) Update(update *UserData) *errors.APIError {
	user, err := u.Service.UpdateUser(u.Id, update)
	if err != nil {
		return err
	}

	*u = user

	return nil
}

func (u *User) Delete() *errors.APIError {
	return u.Service.DeleteUser(u.Id)
}

func (u *User) RefreshAllConnections() *errors.APIError {
	return u.Service.RefreshAllConnections(u.Id)
}

func (u *User) ListAllConnections(filter *utilities.FilterBuilder) (ConnectionList, *errors.APIError) {
	return u.Service.ListAllConnections(u.Id, filter)
}

func (u *User) GetAccount(accountId string) (Account, *errors.APIError) {
	return u.Service.GetAccount(u.Id, accountId)
}

func (u *User) GetAccounts(filter *utilities.FilterBuilder) (AccountsList, *errors.APIError) {
	return u.Service.GetAccounts(u.Id, filter)
}

func (u *User) GetTransaction(transactionId string) (Transaction, *errors.APIError) {
	return u.Service.GetTransaction(u.Id, transactionId)
}

func (u *User) GetTransactions(filter *utilities.FilterBuilder) (TransactionsList, *errors.APIError) {
	return u.Service.GetTransactions(u.Id, filter)
}
