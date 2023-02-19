package pizzahut

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-pizza-api/internal/request"
	"log"
	"strings"
)

const (
	storeURL = "https://api.pizzahut.io/v1/huts?postcode="
	menuURL  = "https://api.pizzahut.io/v1/content/products?sector=uk-1&locale=en-gb"
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
	Delivery   bool    `json:"delivery"`
	Rule       string  `json:"rule"`
	Hidden     bool    `json:"hidden"`
	Locked     bool    `json:"locked"`
}

// getDiscounts will return all of the pizzahut discount codes and vouchers.
func getDiscounts(storeID string, discountChan chan<- []DiscountItem) {

	endpoint := "https://api.pizzahut.io/v2/products/discounts?hutid=" + storeID + "&sector=uk-1&delivery=true"

	body := request.UserAgentGetReq(endpoint)

	discounts := []DiscountItem{}
	err := json.Unmarshal([]byte(body), &discounts)
	if err != nil {
		discountChan <- []DiscountItem{}
		fmt.Print("error: couldn't unmarshal data")
		return
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
func getAllDeals(storeID string, dealsChan chan<- []Deals) {
	endpoint := "https://api.pizzahut.io/v2/products/deals?hutid=" + storeID + "&sector=uk-1&delivery=true"

	body := request.UserAgentGetReq(endpoint)

	deals := []Deals{}
	err := json.Unmarshal([]byte(body), &deals)
	if err != nil {
		dealsChan <- []Deals{}
		return
	}

	dealsChan <- deals
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

func filterAvailableDiscounts(discounts []DiscountItem) []DiscountItem {
	availableDiscount := make([]DiscountItem, 0)

	for _, disc := range discounts {
		if !disc.Hidden && !disc.Locked {
			availableDiscount = append(availableDiscount, disc)
		}
	}

	return availableDiscount
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
// this will return a menu item just to keep data consistent
func lookupDiscountData(discounts []DiscountItem, menu []MenuItem) []MenuItem {
	availableDiscountData := make([]MenuItem, 0)

	// iterate over every deal
	for _, disc := range discounts {
		// iterate over every item in the discounts available
		for _, item := range menu {
			// if the deal is in the discount extract the data that is needed.
			if item.ID == disc.ID { //  && (item.Type == "deal" || item.OtherType == "deal")
				// add required metadata
				item.Price = disc.Discount
				availableDiscountData = append(availableDiscountData, item)
			}
		}
	}

	return availableDiscountData
}

func GetDeals(postcode string) ([]MenuItem, []MenuItem, error) {
	menuChan := make(chan []MenuItem)
	discountChan := make(chan []DiscountItem)
	dealsChan := make(chan []Deals)

	go getMenu(menuChan)

	storeID, err := getStoreID(postcode)
	if err != nil {
		return nil, nil, err
	}

	go getDiscounts(storeID, discountChan)
	go getAllDeals(storeID, dealsChan)

	allDiscounts, allDeals := <-discountChan, <-dealsChan

	availableDeals := filterAvailableDeals(allDeals)
	availableDiscounts := filterAvailableDiscounts(allDiscounts)

	menu := <-menuChan

	dealData := lookupDealData(availableDeals, menu)
	discountData := lookupDiscountData(availableDiscounts, menu)

	return unique(dealData), unique(discountData), nil
}

// https://stackoverflow.com/questions/57706801/deduplicate-array-of-structs
// TODO: Simplify this function
// unique is required because of the pizzahut api being a pile of wank
// need to remove duplicated deals.
func unique(sample []MenuItem) []MenuItem {
	var unique []MenuItem

	type key struct{ value string }

	m := make(map[key]int)

	for _, v := range sample {
		k := key{strings.ToLower(v.ID)}
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
