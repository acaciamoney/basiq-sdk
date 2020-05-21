package errors

import (
	"encoding/json"
	"strings"
)

type ErrorSource struct {
	Pointer   string
	Parameter string
}

type ResponseError struct {
	CorrelationId string `json:"correlationId"`
	Data          []ResponseErrorItem
}

type ResponseErrorItem struct {
	Code   string      `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Source ErrorSource `json:"source"`
}

type APIError struct {
	Data       map[string]interface{}
	Message    string
	Response   ResponseError
	StatusCode int
}

func ParseError(body []byte) (ResponseError, error) {
	var data ResponseError
	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (e *ResponseError) GetMessages() string {
	messages := make([]string, len(e.Data))
	for _, v := range e.Data {
		messages = append(messages, v.Detail)
	}

	return strings.Trim(strings.Join(messages, " "), " ")
}
