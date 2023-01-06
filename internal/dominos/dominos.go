package dominos

import (
	"encoding/json"
	"go-pizza-api/internal/request"
	"strconv"
)

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

func dominoStoreLocator(postcode string) (StoreData, error) {
	// Construct the endpoint URL.
	endpoint := "https://www.dominos.co.uk/storefindermap/storesearch?searchText=" + postcode
	// fmt.Println("dominos store locator endpoint: " + endpoint)

	// Create the client and send a request to the endpoint.
	body := request.DominosGet(endpoint)

	// Populate structs with requests response.
	sd := Stores{}
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
	Desc  string  `json:"description"`
	Id    int     `json:"id"`
	Price float64 `json:"price"`
}

func GetDominosDeals(dealsChan chan<- []DominosStoreDeals, menuID, storeID string) {

	endpoint := "https://www.dominos.co.uk/Deals/StoreDealGroups?dealsVersion=" + menuID + "&fulfilmentMethod=1&isoCode=en-GB&storeId=" + storeID

	body := request.DominosGet(endpoint)

	// fmt.Printf("Dominos Deals Body : %s", string(body))

	sd := []DominosStoreDeals{}
	err := json.Unmarshal([]byte(body), &sd)
	if err != nil {
		dealsChan <- []DominosStoreDeals{}
	}

	// fmt.Printf("Dominos Deals : %+v", sd)

	dealsChan <- sd
}

type Vouchers []struct {
	Desc string `json:"description"`
}

func GetDominosVouchers(voucherChan chan<- Vouchers, menuID, storeID string) {
	endpoint := "https://www.dominos.co.uk/Deals/StoreDealsVouchers?fulfilmentMethod=1&storeId=" + storeID + "&v=120.1.0.8&vouchersOnlineVersion=" + menuID

	body := request.DominosGet(endpoint)

	// fmt.Printf("Dominos Vouchers Body : %s", string(body))

	vouchers := Vouchers{}
	err := json.Unmarshal([]byte(body), &vouchers)
	if err != nil {
		voucherChan <- Vouchers{}
	}

	// fmt.Printf("Dominos Vouchers : %+v", vouchers)

	voucherChan <- vouchers
}

func GetDominosDealsVouchers(postcode string) ([]DominosStoreDeals, Vouchers, error) {
	// fmt.Println("GetDominosDealsVouchers called")
	dealsChan := make(chan []DominosStoreDeals, 1)
	voucherChan := make(chan Vouchers, 1)

	storeData, err := dominoStoreLocator(postcode)
	if err != nil {
		return nil, nil, err
	}

	// fmt.Printf("Dominos Store Data : %+v", storeData)

	go GetDominosDeals(dealsChan, strconv.Itoa(storeData.MenuId), strconv.Itoa(storeData.Id))
	go GetDominosVouchers(voucherChan, strconv.Itoa(storeData.MenuId), strconv.Itoa(storeData.Id))

	return <-dealsChan, <-voucherChan, nil
}
