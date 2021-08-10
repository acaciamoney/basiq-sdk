package v2

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/acaciamoney/basiq-sdk/errors"
)

type Merchant struct {
	BusinessName string `dynamodbav:"businessName"`
	Website      string `dynamodbav:"website"`
	PhoneNumber  struct {
		Local         string `dynamodbav:"local"`
		International string `dynamodbav:"international"`
	} `dynamodbav:"phoneNumber"`
}

type Location struct {
	RouteNo          string `dynamodbav:"routeNo"`
	Route            string `dynamodbav:"route"`
	PostalCode       string `dynamodbav:"postalCode"`
	Suburb           string `dynamodbav:"suburb"`
	State            string `dynamodbav:"state"`
	Country          string `dynamodbav:"country"`
	FormattedAddress string `dynamodbav:"formattedAddress"`
	Geometry         struct {
		Lat string `dynamodbav:"lat"`
		Lng string `dynamodbav:"lng"`
	} `dynamodbav:"geometry"`
}

type Category struct {
	ANZSIC struct {
		Division struct {
			Code  string `dynamodbav:"code"`
			Title string `dynamodbav:"title"`
		} `dynamodbav:"division"`
		Subdivision struct {
			Code  string `dynamodbav:"code"`
			Title string `dynamodbav:"title"`
		} `dynamodbav:"subdivision"`
		Group struct {
			Code  string `dynamodbav:"code"`
			Title string `dynamodbav:"title"`
		} `dynamodbav:"group"`
		Class struct {
			Code  string `dynamodbav:"code"`
			Title string `dynamodbav:"title"`
		} `dynamodbav:"class"`
	} `dynamodbav:"anzsic"`
}

type Enrich struct {
	Type      string `dynamodbav:"type"`
	Class     string `dynamodbav:"class"`
	Direction string `dynamodbav:"direction"`
	Data      struct {
		Merchant Merchant `dynamodbav:"merchant"`
		Location Location `dynamodbav:"location"`
		Category Category `dynamodbav:"category"`
	} `dynamodbav:"data"`
	Links struct {
		Self       string `dynamodbav:"self"`
		LogoMaster string `dynamodbav:"logo-master"`
		LogoThumb  string `dynamodbav:"logo-thumb"`
	} `dynamodbav:"links"`
}

type EnrichService struct {
	Session *Session
}

func NewEnrichService(session *Session) *EnrichService {
	return &EnrichService{
		Session: session,
	}
}

func (es *EnrichService) GetEnrichedTransaction(transaction Transaction) (Enrich, *errors.APIError) {
	var data Enrich
	es.Session.Api.SetHeader("Content-Type", "application/json")
	queryDescription := url.QueryEscape(transaction.Description)
	body, _, err := es.Session.Api.Send("GET", "enrich?q="+queryDescription+"&country=AU&institution="+transaction.Institution, nil)
	if err != nil {
		log.Print("[ERROR] - Failed to make request to enrich service: " + err.Message + "(" + transaction.Id + "|" + transaction.Description + ")")
		return data, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Print("[ERROR] - Failed to parse response from enrich service")
		return data, &errors.APIError{Message: err.Error()}
	}
	return data, nil
}
