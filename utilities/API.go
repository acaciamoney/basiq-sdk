package utilities

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/sethgrid/pester"
)

type API struct {
	host    string
	headers map[string]string
}

func NewAPI(host string) *API {
	return &API{
		host: host,
	}
}

func (api *API) Send(method string, path string, data []byte) ([]byte, int, *errors.APIError) {
	log.Println("Requesting: " + method + "_" + api.host + path)

	var req *http.Request
	var err error

	if data != nil {
		req, err = http.NewRequest(method, api.host+path, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, api.host+path, nil)
	}

	c := pester.New()
	c.MaxRetries = 5
	c.Backoff = pester.ExponentialBackoff
	c.KeepLog = true

	if err != nil {
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	for k, v := range api.headers {
		req.Header.Add(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	if resp.StatusCode > 299 {
		response, err := errors.ParseError(body)
		if err != nil {
			return nil, 0, &errors.APIError{Message: err.Error()}
		}

		return nil, 0, &errors.APIError{
			Response:   response,
			Message:    response.GetMessages(),
			StatusCode: resp.StatusCode,
		}
	}

	return body, resp.StatusCode, nil
}

func (api *API) SetHeader(header string, value string) *API {
	if api.headers == nil {
		api.headers = make(map[string]string)
	}
	api.headers[header] = value

	return api
}
