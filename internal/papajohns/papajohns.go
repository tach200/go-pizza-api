package papajohns

import (
	"encoding/json"
	"go-pizza-api/internal/request"
	"strconv"
)

const (
	storeEndpoint = "https://api2.papajohns.co.uk/api/v1/Store/delivery/"
	dealsEndpoint = "https://api2.papajohns.co.uk/api/v1/deal/"
)

type StoreID struct {
	ID int `json:"id"`
}

type StoreInfo struct {
	Version int     `json:"version"`
	Data    StoreID `json:"data"`
}

// getStoreInfo returns information about the store closest to that postcode
// this information can be used to populate other endpoint information.
func getStoreInfo(postcode string) (StoreInfo, error) {
	endpoint := storeEndpoint + postcode

	body := request.UserAgentGetReq(endpoint)

	storeData := StoreInfo{}
	err := json.Unmarshal([]byte(body), &storeData)
	if err != nil {
		return storeData, err
	}

	return storeData, nil
}

type Deal struct {
	DisplayName string `json:"displayName"`
	PromoURL    string `json:"promo"`
	Desc        string `json:"description"`
	Available   int    `json:"availibility"`
}

type Deals struct {
	Version int    `json:"version"`
	Deals   []Deal `json:"data"`
}

// getDeals returns the deals that are available in the given store.
func GetDeals(postcode string) (Deals, error) {
	storeID, err := getStoreInfo(postcode)
	if err != nil {
		return Deals{}, err
	}

	endpoint := dealsEndpoint + strconv.Itoa(storeID.Data.ID)

	body := request.UserAgentGetReq(endpoint)

	var deals Deals
	err = json.Unmarshal([]byte(body), &deals)
	if err != nil {
		return deals, err
	}

	return deals, nil
}
