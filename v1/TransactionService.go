package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/acaciamoney/basiq-sdk/utilities"
)

type TransactionsList struct {
	Count   int               `json:"count"`
	Data    []Transaction     `json:"data"`
	Links   map[string]string `json:"links"`
	Service *TransactionService
}

type Transaction struct {
	Id              string                 `json:"id"`
	Status          string                 `json:"status"`
	Description     string                 `json:"description"`
	PostDate        string                 `json:"postDate"`
	TransactionDate string                 `json:"transactionDate"`
	Amount          string                 `json:"amount"`
	Balance         string                 `json:"balance"`
	BankCategory    string                 `json:"bankCategory"`
	Account         string                 `json:"account"`
	Institution     string                 `json:"institution"`
	Connection      string                 `json:"connection"`
	Class           map[string]interface{} `json:"class"`
}

type TransactionService struct {
	Session *Session
}

func NewTransactionService(session *Session) *TransactionService {
	return &TransactionService{
		Session: session,
	}
}

func (ts *TransactionService) GetTransactions(userId string, filter *utilities.FilterBuilder) (TransactionsList, *errors.APIError) {
	var data TransactionsList

	url := "users/" + userId + "/transactions"

	if filter != nil {
		url = url + "?" + filter.GetFilter()
	}
	body, _, err := ts.Session.Api.Send("GET", url, nil)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return data, &errors.APIError{Message: err.Error()}
	}

	data.Service = ts

	return data, nil
}

func (tl *TransactionsList) Next() (bool, *errors.APIError) {
	var data TransactionsList

	if next, ok := tl.Links["next"]; ok {
		nextPath := next[strings.LastIndex(next, ".io/")+4:]
		body, _, err := tl.Service.Session.Api.Send("GET", nextPath, nil)
		if err != nil {
			return false, err
		}

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println(string(body))
			return false, &errors.APIError{Message: err.Error()}
		}

		data.Service = tl.Service
		*tl = data

		return true, nil
	}

	return false, nil
}
