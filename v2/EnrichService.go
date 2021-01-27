package v2

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/acaciamoney/basiq-sdk/errors"
)

type Merchant struct {
	BusinessName string `json:"businessName"`
	Website      string `json:"website"`
	PhoneNumber  struct {
		Local         string `json:"local"`
		International string `json:"international"`
	} `json:"phoneNumber"`
}

type Location struct {
	RouteNo          string `json:"routeNo"`
	Route            string `json:"route"`
	PostalCode       string `json:"postalCode"`
	Suburb           string `json:"suburb"`
	State            string `json:"state"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formattedAddress"`
	Geometry         struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	} `json:"geometry"`
}

type Category struct {
	ANZSIC struct {
		Division struct {
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"division"`
		Subdivision struct {
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"subdivision"`
		Group struct {
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"group"`
		Class struct {
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"class"`
	} `json:"anzsic"`
}

type Enrich struct {
	Type      string `json:"type"`
	Class     string `json:"class"`
	Direction string `json:"direction"`
	Data      struct {
		Merchant Merchant `json:"merchant"`
		Location Location `json:"location"`
		Category Category `json:"category"`
	} `json:"data"`
	Links struct {
		Self       string `json:"self"`
		LogoMaster string `json:"logoMaster"`
		LogoThumb  string `json:"logoThumb"`
	} `json:"links"`
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
		log.Print("[ERROR] - Failed to make request to enrich service")
		return data, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Print("[ERROR] - Failed to parse response from enrich service")
		return data, &errors.APIError{Message: err.Error()}
	}
	return data, nil
}
