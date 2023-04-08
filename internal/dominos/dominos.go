package dominos

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-pizza-api/internal/ranking"
	"go-pizza-api/internal/request"
	"regexp"
	"strconv"
	"strings"
)

const percentRegex = "(\\d+(\\.\\d+)?%)"

// Structs Hold information about the store
type StoreData struct {
	Id         int  `json:"id"`
	Open       bool `json:"isOpen"`
	CanDeliver bool `json:"localStoreCanDeliverToAddress"`
	MenuId     int  `json:"MenuVersion"`
}
type Stores struct {
	Store StoreData `json:"localStore"`
}

func getStoreID(postcode string) (StoreData, error) {
	endpoint := "https://www.dominos.co.uk/storefindermap/storesearch?searchText=" + postcode

	body := request.Get(endpoint)

	sd := Stores{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		return sd.Store, err
	}

	return sd.Store, nil
}

type StoreDeals struct {
	StoreDeals []Deals `json:"storeDeals"`
}
type Deals struct {
	Deal []Deal `json:"deals"`
	Name string `json:"name"`
}

type Deal struct {
	Desc        string        `json:"description"`
	Id          int           `json:"id"`
	Price       float64       `json:"price"`
	DealContent []DealContent `json:"steps"`
}

type DealContent struct {
	Product string `json:"imageUrl"`
}

func getDeals(dealsChan chan<- []StoreDeals, menuID, storeID string) {

	endpoint := "https://www.dominos.co.uk/Deals/StoreDealGroups?dealsVersion=" + menuID + "&fulfilmentMethod=1&isoCode=en-GB&storeId=" + storeID

	body := request.Get(endpoint)

	fmt.Print(string(body))

	sd := []StoreDeals{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		dealsChan <- []StoreDeals{}
	}

	dealsChan <- sd
}

type Vouchers []struct {
	Desc     string   `json:"description"`
	MinSpend MinSpend `json:"minimumSpend"`
}

type MinSpend struct {
	Amount float64 `json:"amount"`
}

func getVouchers(voucherChan chan<- Vouchers, menuID, storeID string) {
	endpoint := "https://www.dominos.co.uk/Deals/StoreDealsVouchers?fulfilmentMethod=1&storeId=" + storeID + "&v=120.1.0.8&vouchersOnlineVersion=" + menuID

	body := request.Get(endpoint)

	vouchers := Vouchers{}
	err := json.Unmarshal([]byte(body), &vouchers)
	if err != nil {
		voucherChan <- Vouchers{}
	}

	voucherChan <- vouchers
}

// GetAllSavings gets deals and vouchers from dominos
func GetAllSavings(postcode string) ([]StoreDeals, Vouchers, error) {
	dealsChan := make(chan []StoreDeals, 1)
	voucherChan := make(chan Vouchers, 1)

	storeData, err := getStoreID(postcode)
	if err != nil {
		return nil, nil, err
	}

	go getDeals(dealsChan, strconv.Itoa(storeData.MenuId), strconv.Itoa(storeData.Id))
	go getVouchers(voucherChan, strconv.Itoa(storeData.MenuId), strconv.Itoa(storeData.Id))

	return <-dealsChan, <-voucherChan, nil
}

// GetReduction uses a regular expression to extract the reduction amount from the deal
// this is because dominos doesn't expose that data in the api in a conveinent format
// because the data is extracted from a string it also needs to be converted to a float
// for easier arithmetic.
func GetReduction(desc string) (float64, error) {
	regx := regexp.MustCompile(percentRegex)

	reductionStr := regx.FindString(desc)

	if reductionStr == "" {
		return -1, errors.New("no keywords extracted from text")
	}

	// ready the string for conversion to float
	reductionStr = strings.Trim(reductionStr, "%")

	reduction, err := strconv.ParseFloat(reductionStr, 64)
	if err != nil {
		return -1, err
	}

	return reduction, nil
}

func FormatProductData(dealContent []DealContent, dealDesc string) []ranking.Product {
	productData := make([]ranking.Product, 0)
	prevProductName := ""

	product := ranking.Product{}
	prodCount := 0

	for _, d := range dealContent {

		prodType := d.Product[27 : len(d.Product)-4]
		if prodType == "pizza" {
			prodType = strings.ToLower(ranking.GetPizzaSize(dealDesc)) + " pizza"
		}

		if prevProductName != prodType && prevProductName != "" {
			productData = append(productData, product)
			prodCount = 0
		}

		prodCount++
		product.ProductType = prodType
		product.ProductCount = prodCount
		prevProductName = prodType
	}

	productData = append(productData, product)

	return productData
}
