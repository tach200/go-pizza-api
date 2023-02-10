package pizzahut

import (
	"encoding/json"
	"errors"
	"go-pizza-api/internal/request"
	"log"
	"strings"
)

const (
	storeURL     = "https://api.pizzahut.io/v1/huts?postcode="
	menuURL      = "https://api.pizzahut.io/v1/content/products?sector=uk-1&locale=en-gb"
	discountsURL = "https://api.pizzahut.io/v1/content/products?sector=uk-1&locale=en-gb"
)

type Store struct {
	ID string `json:"id"`
}

// getStoreDetails returns the ID of the store closest to the postcoce given
// if there is no store available then assume that one cannot deliver.
func getStoreID(postcode string) (string, error) {
	endpoint := storeURL + postcode

	body := request.UserAgentGetReq(endpoint)

	store := []Store{} // needs to be an array to unmarshal
	err := json.Unmarshal([]byte(body), &store)
	if err != nil {
		log.Fatal(err)
		return store[0].ID, err
	}

	if store[0].ID == "" {
		return "", errors.New("error : delivery not available for this postcode")
	}

	return store[0].ID, nil
}

type MenuItem struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Desc      string  `json:"desc"`
	Type      string  `json:"productType"`
	OtherType string  `json:"type"`
	Price     float64 `json:"price"`
}

// getMenu will return all of the pizzahut menu items.
func getMenu(menuChan chan<- []MenuItem) {
	body := request.UserAgentGetReq(menuURL)

	menu := []MenuItem{}
	err := json.Unmarshal([]byte(body), &menu)
	if err != nil {
		log.Fatal(err)
	}

	menuChan <- menu
}

type DiscountItem struct {
	ID         string  `json:"id"`
	Discount   float64 `json:"amount"`
	Collection bool    `json:"collection"`
	Delivery   bool    `json:"threshold"`
	Rule       string  `json:"rule"`
}

// getDiscounts will return all of the pizzahut discount codes and vouchers.
func getDiscounts(discountChan chan<- []DiscountItem) {
	body := request.UserAgentGetReq(discountsURL)

	// Put the data into a struct.
	discounts := []DiscountItem{}
	err := json.Unmarshal([]byte(body), &discounts)
	if err != nil {
		log.Fatal(err)
	}

	// Get the menu ready for elsewhere.
	discountChan <- discounts
}

type Deals struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Hidden   bool    `json:"hidden"`
	Locked   bool    `json:"locked"`
	Delivery bool    `json:"delivery"`
}

// getAllDeals will extract all the deals that the API returns for the given storeID
// however not all of these deals will be available to the user so some further filtering is requited.
func getAllDeals(storeID string) ([]Deals, error) {
	endpoint := "https://api.pizzahut.io/v2/products/deals?hutid=" + storeID + "&sector=uk-1&delivery=true"

	body := request.UserAgentGetReq(endpoint)

	deals := []Deals{}
	err := json.Unmarshal([]byte(body), &deals)
	if err != nil {
		return nil, err
	}

	return deals, nil
}

// filterAvailableDeals will return the deals that are available to the customer.
func filterAvailableDeals(deals []Deals) []Deals {
	availableDeals := make([]Deals, 0)

	for _, deal := range deals {
		if !deal.Hidden && !deal.Locked {
			availableDeals = append(availableDeals, deal)
		}
	}
	return availableDeals
}

// lookupDealData will find the deal description, price and other important metadata
// it will return menu items are these are more detailed
// this sucks but the pizzahut api has forced my hand.
// this probably be optimised, but for now fuck it.
func lookupDealData(deals []Deals, menu []MenuItem) []MenuItem {
	availableDealData := make([]MenuItem, 0)

	// iterate over every deal
	for _, deal := range deals {
		// iterate over every item in the menu
		for _, item := range menu {
			// if the deal is in the menu extract the data that is needed.
			if item.ID == deal.ID { //  && (item.Type == "deal" || item.OtherType == "deal")
				// add required metadata
				item.Price = deal.Price
				availableDealData = append(availableDealData, item)
			}
		}
	}

	return availableDealData
}

// lookupDiscountData will find the discount description, price and other important metadata
// func lookupDiscountData(deals []Deals, discounts []DiscountItem) []MenuItem {
// 	return nil
// }

func GetDeals(postcode string) ([]MenuItem, error) {
	menuChan := make(chan []MenuItem)
	// discountChan := make(chan []DiscountItem)

	// fetch all of the items and discount codes.
	go getMenu(menuChan)
	// go getDiscounts(discountChan)

	storeID, err := getStoreID(postcode)
	if err != nil {
		return nil, err
	}

	allDeals, err := getAllDeals(storeID)
	if err != nil {
		return nil, err
	}

	availableDeals := filterAvailableDeals(allDeals)

	dealData := lookupDealData(availableDeals, <-menuChan)
	// discountData := lookupDiscountData(availableDeals, <-discountChan)

	return dealData, nil
}

// https://stackoverflow.com/questions/57706801/deduplicate-array-of-structs
func unique(sample []MenuItem) []MenuItem {
	var unique []MenuItem

	type key struct{ value1, value2 string }

	m := make(map[key]int)

	for _, v := range sample {
		k := key{strings.ToLower(v.Desc), v.ID}
		if i, ok := m[k]; ok {
			// Overwrite previous value per requirement in
			// question to keep last matching value.
			unique[i] = v
		} else {
			// Unique key found. Record position and collect
			// in result.
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}

// TODO:
// isStudent
// Discounts
// NHS
// Collection and Delivery is seperate API call like wtf...
