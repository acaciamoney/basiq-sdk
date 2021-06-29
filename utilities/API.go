package utilities

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/acaciamoney/basiq-sdk/errors"
	"github.com/sethgrid/pester"
)

type API struct {
	host    string
	headers map[string]string
	mutex   sync.Mutex
}

func NewAPI(host string) *API {
	return &API{
		host:  host,
		mutex: sync.Mutex{},
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
	c.Concurrency = 1
	c.MaxRetries = 20
	c.Backoff = pester.ExponentialJitterBackoff
	c.KeepLog = true

	if err != nil {
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	api.mutex.Lock()
	for k, v := range api.headers {
		req.Header.Add(k, v)
	}
	api.mutex.Unlock()

	resp, err := c.Do(req)

	if err != nil {
		log.Print("[ERROR] - Unable to send request to basiq API")
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("[ERROR] - Unable to parse response from basiq API")
		return nil, 0, &errors.APIError{Message: err.Error()}
	}

	if resp.StatusCode > 299 {
		response, err := errors.ParseError(body)
		if err != nil {
			log.Print("[ERROR] - Unable to parse error from basiq API")
			return nil, 0, &errors.APIError{Message: err.Error()}
		}
		log.Print("[ERROR] - Bad response code from basiq API")
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
