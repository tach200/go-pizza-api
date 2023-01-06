package papajohns

import (
	"encoding/json"
	"go-pizza-api/internal/request"
	"strconv"
	"time"
)

const (
	storeEndpoint = "https://api2.papajohns.co.uk/api/v1/Store/delivery/"
	dealsEndpoint = "https://api2.papajohns.co.uk/api/v1/deal/"
)

var (
	daysOfWeek = map[int]string{
		1: "Monday",
		2: "Tuesday",
		3: "Wednesday",
		4: "Thursday",
		5: "Friday",
		6: "Saturday",
		7: "Sunday",
	}
)

type StoreID struct {
	ID int `json:"id"`
}

type StoreInfo struct {
	Data StoreID `json:"data"`
}

// getStoreInfo returns information about the store closest to that postcode
// this information can be used to populate other endpoint information.
func getStoreInfo(postcode string) (StoreInfo, error) {
	endpoint := storeEndpoint + postcode

	body := request.PapaGet(endpoint)

	storeData := StoreInfo{}
	err := json.Unmarshal([]byte(body), &storeData)
	if err != nil {
		return storeData, err
	}

	return storeData, nil
}

type Deal struct {
	DisplayName string   `json:"name"`
	PromoURL    string   `json:"promo"`
	Desc        string   `json:"description"`
	Displayed   bool     `json:"showOnDealsPage"`
	Available   int      `json:"availability"`
	Price       float64  `json:"price"`
	Schedule    Schedule `json:"schedule"`
	StudentDeal bool     `json:"studentDeal"`
}

type Schedule struct {
	DaysOfWeek []int `json:"daysOfWeek"`
}

type Deals struct {
	Deals []Deal `json:"data"`
}

// getDeals returns the deals that are available in the given store.
func GetDeals(postcode string) ([]Deal, error) {
	storeID, err := getStoreInfo(postcode)
	if err != nil {
		return []Deal{}, err
	}

	endpoint := dealsEndpoint + strconv.Itoa(storeID.Data.ID)

	body := request.PapaGet(endpoint)

	var deals Deals
	err = json.Unmarshal([]byte(body), &deals)
	if err != nil {
		return deals.Deals, err
	}

	return scheduleFilter(deals.Deals), nil
}

// scheduleFilter will remove deals that are not available
func scheduleFilter(allDeals []Deal) []Deal {
	var availableDeals []Deal

	// fmt.Print(allDeals)

	for _, deal := range allDeals {
		if len(deal.Schedule.DaysOfWeek) == 0 {
			availableDeals = append(availableDeals, deal)
			continue
		}

		for _, day := range deal.Schedule.DaysOfWeek {
			if daysOfWeek[day] == time.Now().Weekday().String() {
				availableDeals = append(availableDeals, deal)
				break
			}
		}
	}

	return availableDeals
}
