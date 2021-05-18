package v1

type AccountsList struct {
	Count int       `json:"count"`
	Data  []Account `json:"data"`
}

type Account struct {
	Id             string                 `json:"id"`
	AccountNo      string                 `json:"accountNo"`
	Name           string                 `json:"name"`
	Currency       string                 `json:"currency"`
	Balance        string                 `json:"balance"`
	AvailableFunds string                 `json:"availableFunds"`
	LastUpdated    string                 `json:"lastUpdated"`
	Class          map[string]interface{} `json:"class"`
	Status         string                 `json:"status"`
	Institution    string                 `json:"institution"`
	Connection     string                 `json:"connection"`
}
