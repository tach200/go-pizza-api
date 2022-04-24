package dominos

import (
	"encoding/json"
	"fmt"
	"go-pizza-api/internal/request"
	"strconv"
)

// Structs Hold information about the store
type DominosStoreDetails struct {
	Id         int  `json:"id"`
	Open       bool `json:"isOpen"`
	CanDeliver bool `json:"localStoreCanDeliverToAddress"`
	MenuId     int  `json:"MenuVersion"`
}
type DominosStore struct {
	Store DominosStoreDetails `json:"localStore"`
}

func dominoStoreLocator(postcode string) (DominosStoreDetails, error) {
	// Construct the endpoint URL.
	endpoint := "https://www.dominos.co.uk/storefindermap/storesearch?searchText=" + postcode
	// fmt.Println("dominos store locator endpoint: " + endpoint)

	// Create the client and send a request to the endpoint.
	body := request.DominosGet(endpoint)

	// Populate structs with requests response.
	sd := DominosStore{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		return sd.Store, err
	}

	return sd.Store, nil
}

// Dominos deals data, data arrives in horrible format.
// Structs are built for the JSON data
type DominosStoreDeals struct {
	StoreDeals []DominosDeals `json:"storeDeals"`
}
type DominosDeals struct {
	Deal []DominosDeal `json:"deals"`
	Name string        `json:"name"`
}

type DominosDeal struct {
	Desc string `json:"description"`
	Id   int    `json:"id"`
}

func GetDominosDeals(postcode string) ([]DominosStoreDeals, error) {
	// 1. Retrieve dominos store data to construct the endpoint URL.
	storeData, err := dominoStoreLocator(postcode)
	if err != nil {
		return nil, err
	}
	endpoint := "https://www.dominos.co.uk/Deals/StoreDealGroups?dealsVersion=" + strconv.Itoa(storeData.MenuId) + "&fulfilmentMethod=1&isoCode=en-GB&storeId=" + strconv.Itoa(storeData.Id)
	fmt.Println("dominos deal endpoint: " + endpoint)

	// 2. Make a request to the endpoint.
	body := request.DominosGet(endpoint)

	// 4. Populate structs with requests response.
	sd := []DominosStoreDeals{}
	err = json.Unmarshal([]byte(body), &sd)
	if err != nil {
		return nil, err
	}

	return sd, nil
}
