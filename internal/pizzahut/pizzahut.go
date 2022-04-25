package pizzahut

import (
	"encoding/json"
	"go-pizza-api/internal/request"
	"log"
)

type PizzahutDetails struct {
	StoreID string `json:"id"`
}

func pizzahutStoreLocator(postcode string) (string, error) {
	// Construct URL endpoint.
	endpoint := "https://api.pizzahut.io/v1/huts?postcode=" + postcode

	// Make a request to the endpoint.
	body := request.UserAgentGetReq(endpoint)

	// Put the JSON data into a struct.
	sd := []PizzahutDetails{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		log.Fatal(err)
		return sd[0].StoreID, err
	}
	// The first store is the one available for delivery.
	return sd[0].StoreID, nil
}

type PizzahutMenuItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func getPizzahutMenu(menu chan<- []PizzahutMenuItem) {
	// Make ednpoint
	endpoint := "https://api.pizzahut.io/v1/content/products?sector=uk-1&locale=en-gb"
	// Make a request to the Menu
	body := request.UserAgentGetReq(endpoint)

	// Put the data into a struct.
	sd := []PizzahutMenuItem{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		log.Fatal(err)
	}

	// Get the menu ready for elsewhere.
	menu <- sd
}

type PizzahutDeals struct {
	Id     string  `json:"id"`
	Price  float32 `json:"price"`
	Hidden bool    `json:"hidden"`
}

func GetPizzahutDeals(postcode string) ([]PizzahutMenuItem, error) {
	menu := make(chan []PizzahutMenuItem)
	// Make a request to get the Menu endpoint. This will be used in conjuction with the other request.
	go getPizzahutMenu(menu)
	// Construct the endpoint.
	storeData, err := pizzahutStoreLocator(postcode)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	endpoint := "https://api.pizzahut.io/v2/products/deals?hutid=" + storeData + "&sector=uk-1&delivery=true"

	// Make a request to the endpoint. This endpoint will receive ID's for menu items.
	body := request.UserAgentGetReq(endpoint)

	// Put the data into a struct.
	uSd := []PizzahutDeals{}
	err = json.Unmarshal([]byte(body), &uSd)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// API filtering is limited. So have to do some cleaning here.
	var sd []PizzahutDeals
	for _, v := range uSd {
		// Some deals that are not available are still in the json dump.
		// Filter out the hidden deals.
		if !v.Hidden {
			sd = append(sd, v)
		}
	}

	// Use the Filtered ID's to retrieve the menu items. Wait for allItems to arrive.
	var availableItems []PizzahutMenuItem
	for _, v := range <-menu {
		for _, v2 := range sd {
			if v.Id == v2.Id {
				availableItems = append(availableItems, v)
				break
			}
		}
	}
	return availableItems, nil
}
