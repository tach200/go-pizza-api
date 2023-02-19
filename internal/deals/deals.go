package deals

import (
	"go-pizza-api/internal/dominos"
	"strconv"
	"sync"
)

type AllDeals struct {
	Restaurant          string
	DealName            string
	DealDesc            string
	DealUrl             string
	StudentDeal         bool
	CollectionOnly      bool
	Score               float64
	DealPrice           float64
	DeliveryCost        float64
	DealType            string
	Discount            float64
	PriceBeforeDiscount float64
	PriceAfterDiscount  float64
	DealContent         interface{}
}

type Product struct {
	Name  string
	Count int
}

func GetDeals(postcode string) []AllDeals {
	//Create list of structs to store clean data.
	deals := []AllDeals{}
	var wg sync.WaitGroup
	// 3 requests need to be made
	wg.Add(1)

	// Pizzahut
	// go func() {
	// 	defer wg.Done()
	// 	pizzas, discounts, err := pizzahut.GetDeals(postcode)
	// 	if err != nil {
	// 		return
	// 	}
	// 	for _, item := range pizzas {
	// 		deals = append(deals, AllDeals{
	// 			Restaurant:  "Pizza Hut",
	// 			DealName:    item.Title,
	// 			DealDesc:    item.Desc,
	// 			DealPrice:   item.Price,
	// 			DealUrl:     "https://www.pizzahut.co.uk/order/deal/?id=" + item.ID,
	// 			DealType:    item.Type,
	// 			DealContent: item.DealContent,
	// 			Score:       rankDealScore(item.Desc, item.Price, pizzahutSizes),
	// 		})
	// 	}
	// 	for _, disc := range discounts {

	// 		score, costAfterDiscount := scoreDiscount(disc.Reduction, disc.Price)

	// 		deals = append(deals, AllDeals{
	// 			Restaurant:         "Pizza Hut",
	// 			DealName:           disc.Title,
	// 			DealDesc:           disc.Desc,
	// 			DealPrice:          disc.Price,
	// 			DealUrl:            "https://www.pizzahut.co.uk/order/deal/?id=" + disc.ID,
	// 			DealType:           disc.Type,
	// 			PriceAfterDiscount: costAfterDiscount,
	// 			Score:              score,
	// 		})
	// 	}
	// }()

	// // Papajohns
	// go func() {
	// 	defer wg.Done()
	// 	papajohns, err := papajohns.GetDeals(postcode)
	// 	if err != nil {
	// 		return
	// 	}
	// 	for _, item := range papajohns {
	// 		if item.Displayed {

	// 			var collection bool
	// 			if item.ShippingMethod == 1 {
	// 				collection = true
	// 			}

	// 			var score float64
	// 			var costAfterDiscount = -1.00
	// 			if item.ReductionType == "%" {
	// 				score, costAfterDiscount = scoreDiscount(item.Reduction, item.MinimumSpend)
	// 			} else {
	// 				score = rankDealScore(item.Desc, item.Price, papajohsSizes)
	// 			}

	// 			deals = append(deals, AllDeals{
	// 				Restaurant:          "Papa Johns",
	// 				DealName:            item.DisplayName,
	// 				DealUrl:             "https://www.papajohns.co.uk/deals",
	// 				DealDesc:            item.Desc,
	// 				DealPrice:           item.Price,
	// 				DealType:            item.ReductionType,
	// 				StudentDeal:         item.StudentDeal,
	// 				CollectionOnly:      collection,
	// 				Discount:            item.Reduction,
	// 				PriceAfterDiscount:  float64(costAfterDiscount),
	// 				PriceBeforeDiscount: item.MinimumSpend, //field only present on % deals
	// 				DealContent:         item.DealContent,
	// 				Score:               score,
	// 			})
	// 		}
	// 	}
	// }()

	// // Dominos
	go func() {
		defer wg.Done()
		domnios, vouchers, err := dominos.GetAllSavings(postcode)
		if err != nil {
			return
		}
		for _, item := range domnios[0].StoreDeals {
			deals = append(deals, AllDeals{
				Restaurant:     "Domino's",
				DealName:       item.Name,
				DealUrl:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
				DealPrice:      item.Deal[0].Price,
				DealDesc:       item.Deal[0].Desc,
				DealType:       "",
				CollectionOnly: false,
				StudentDeal:    false,
				DealContent:    item.Deal[0].DealContent,
				Score:          rankDealScore(item.Deal[0].Desc, item.Deal[0].Price, dominosSizes),
			})
		}
		for _, item := range vouchers {

			reduction, err := dominos.GetReduction(item.Desc)
			if err != nil {
				return
			}

			score, costAfterDiscount := scoreDiscount(reduction, item.MinSpend.Amount)

			deals = append(deals, AllDeals{
				Restaurant:          "Domino's",
				DealName:            "Savings Voucher",
				DealUrl:             "https://www.dominos.co.uk/deals",
				DealDesc:            item.Desc,
				DealType:            "%",
				CollectionOnly:      false,
				StudentDeal:         false,
				PriceBeforeDiscount: item.MinSpend.Amount,
				PriceAfterDiscount:  costAfterDiscount,
				Score:               score,
			})
		}
	}()

	wg.Wait()
	return deals
}
