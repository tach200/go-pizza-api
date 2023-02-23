package deals

import (
	"go-pizza-api/internal/pizzahut"
	"sync"
)

type AllDeals struct {
	// Data
	Restaurant   string
	DealName     string
	DealDesc     string
	DealUrl      string
	DealPrice    float64
	DeliveryCost float64
	DealType     string

	// Filters
	StudentDeal    bool
	CollectionOnly bool

	Score float64

	// Discounts
	Discount            float64
	PriceBeforeDiscount float64
	PriceAfterDiscount  float64

	// Products
	InchesOfPizza float64
	Products      interface{}
}

func GetDeals(postcode string) []AllDeals {
	//Create list of structs to store clean data.
	deals := []AllDeals{}
	var wg sync.WaitGroup
	// 3 requests need to be made
	wg.Add(1)

	// Pizzahut
	go func() {
		defer wg.Done()
		pizzas, discounts, err := pizzahut.GetDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range pizzas {

			products := pizzahut.FormatProductData(item.DealContent)
			inchesOfPizza := calculateTotalInches(Product(products[0]), pizzahutSizes)

			deals = append(deals, AllDeals{
				Restaurant: "Pizza Hut",
				DealName:   item.Title,
				DealDesc:   item.Desc,
				DealPrice:  item.Price,
				DealUrl:    "https://www.pizzahut.co.uk/order/deal/?id=" + item.ID,
				DealType:   item.Type,
				Products:   products,
				Score:      scoreDeal(inchesOfPizza, item.Price),
			})
		}
		for _, disc := range discounts {

			score, costAfterDiscount := scoreDiscount(disc.Reduction, disc.Price)

			deals = append(deals, AllDeals{
				Restaurant:         "Pizza Hut",
				DealName:           disc.Title,
				DealDesc:           disc.Desc,
				DealPrice:          disc.Price,
				DealUrl:            "https://www.pizzahut.co.uk/order/deal/?id=" + disc.ID,
				DealType:           "%",
				PriceAfterDiscount: costAfterDiscount,
				Score:              score,
			})
		}
	}()

	// Papajohns
	// go func() {
	// 	defer wg.Done()
	// 	papaj, err := papajohns.GetDeals(postcode)
	// 	if err != nil {
	// 		return
	// 	}
	// 	for _, item := range papaj {
	// 		if item.Displayed {

	// 			var collection bool
	// 			// means collection
	// 			if item.ShippingMethod == 1 {
	// 				collection = true
	// 			}

	// 			var products []papajohns.Product
	// 			var score float64
	// 			var costAfterDiscount = -1.00

	// 			if item.ReductionType == "%" {
	// 				score, costAfterDiscount = scoreDiscount(item.Reduction, item.MinimumSpend)
	// 			} else {
	// 				products = papajohns.FormatProductData(item.DealContent)
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
	// 				Products:            products,
	// 				Score:               score,
	// 			})
	// 		}
	// 	}
	// }()

	// // Dominos
	// go func() {
	// 	defer wg.Done()
	// 	domnios, vouchers, err := dominos.GetAllSavings(postcode)
	// 	if err != nil {
	// 		return
	// 	}

	// 	for _, item := range domnios[0].StoreDeals {
	// 		deals = append(deals, AllDeals{
	// 			Restaurant:     "Domino's",
	// 			DealName:       item.Name,
	// 			DealUrl:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
	// 			DealPrice:      item.Deal[0].Price,
	// 			DealDesc:       item.Deal[0].Desc,
	// 			DealType:       "",
	// 			CollectionOnly: false,
	// 			StudentDeal:    false,
	// 			Products:       dominos.FormatProductData(item.Deal[0].DealContent),
	// 			Score:          rankDealScore(item.Deal[0].Desc, item.Deal[0].Price, dominosSizes),
	// 		})
	// 	}
	// 	for _, item := range vouchers {

	// 		reduction, err := dominos.GetReduction(item.Desc)
	// 		if err != nil {
	// 			return
	// 		}

	// 		score, costAfterDiscount := scoreDiscount(reduction, item.MinSpend.Amount)

	// 		deals = append(deals, AllDeals{
	// 			Restaurant:          "Domino's",
	// 			DealName:            "Savings Voucher",
	// 			DealUrl:             "https://www.dominos.co.uk/deals",
	// 			DealDesc:            item.Desc,
	// 			DealType:            "%",
	// 			CollectionOnly:      false,
	// 			StudentDeal:         false,
	// 			PriceBeforeDiscount: item.MinSpend.Amount,
	// 			PriceAfterDiscount:  costAfterDiscount,
	// 			Score:               score,
	// 		})
	// 	}
	// }()

	wg.Wait()
	return deals
}
