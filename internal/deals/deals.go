package deals

import (
	"go-pizza-api/internal/papajohns"
	"sync"
)

type AllDeals struct {
	Restaurant     string  `json:"Restaurant"`
	DealName       string  `json:"DealName"`
	DealDesc       string  `json:"DealDesc"`
	Url            string  `json:"Url"`
	Rank           float64 `json:"Rank"`
	Price          float64 `json:"Price"`
	DealType       string  `json:"DealType"`
	StudentDeal    bool    `json:"StudentDeal"`
	CollectionOnly bool    `json:"CollectionOnly"`
	Discount       float64 `json:"Discount"`
	DeliveryCost   float64 `json:"DeliveryCost"`
}

func GetDeals(postcode string) []AllDeals {
	// fmt.Printf("GetDeals called with postcode: %s", postcode)
	//Create list of structs to store clean data.
	deals := []AllDeals{}
	var wg sync.WaitGroup
	// 3 requests need to be made
	wg.Add(1)

	// go func() {
	// 	defer wg.Done()
	// 	pizzahut, err := pizzahut.GetDeals(postcode)
	// 	if err != nil {
	// 		return
	// 	}
	// 	for _, item := range pizzahut {
	// 		deals = append(deals, AllDeals{
	// 			Restaurant: "Pizza Hut",
	// 			DealName:   item.Title,
	// 			DealDesc:   item.Desc,
	// 			Price:      item.Price,
	// 			Url:        "https://www.pizzahut.co.uk/order/deal/?id=" + item.ID,
	// 			DealType:   item.Type,
	// 			// Rank:       rankScore(item.Desc, item.Price, pizzahutSizes),
	// 		})
	// 	}
	// }()

	go func() {
		defer wg.Done()
		papajohns, err := papajohns.GetDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range papajohns {
			if item.Displayed {

				var collection bool

				if item.ShippingMethod == 1 {
					collection = true
				}

				deals = append(deals, AllDeals{
					Restaurant:     "Papa Johns",
					DealName:       item.DisplayName,
					Url:            "https://www.papajohns.co.uk/deals",
					DealDesc:       item.Desc,
					Price:          item.Price,
					DealType:       item.ReductionType,
					StudentDeal:    item.StudentDeal,
					CollectionOnly: collection,
					Discount:       item.Reduction,
					Rank:           rankScore(item.Desc, item.Price, papajohsSizes),
				})
			}
		}
	}()

	// go func() {
	// 	defer wg.Done()
	// 	domDeals, vouchers, err := dominos.GetDominosDealsVouchers(postcode)
	// 	if err != nil {
	// 		return
	// 	}
	// 	for _, item := range domDeals[0].StoreDeals {
	// 		deals = append(deals, AllDeals{
	// 			Restaurant:     "Domino's",
	// 			DealName:       item.Name,
	// 			Url:            "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
	// 			Price:          item.Deal[0].Price,
	// 			DealDesc:       item.Deal[0].Desc,
	// 			DealType:       "",
	// 			CollectionOnly: false,
	// 			StudentDeal:    false,
	// 			Rank:           rankScore(item.Deal[0].Desc, item.Deal[0].Price, dominosSizes),
	// 		})
	// 	}
	// 	for _, item := range vouchers {
	// 		deals = append(deals, AllDeals{
	// 			Restaurant:     "Domino's",
	// 			DealName:       "Savings Voucher",
	// 			Url:            "https://www.dominos.co.uk/deals",
	// 			DealDesc:       item.Desc,
	// 			DealType:       "%",
	// 			CollectionOnly: false,
	// 			StudentDeal:    false,
	// 			Rank:           rankScore(item.Desc, 0, dominosSizes),
	// 		})
	// 	}
	// }()

	wg.Wait()
	return deals
}
