package deals

import (
	"go-pizza-api/internal/dominos"
	"go-pizza-api/internal/papajohns"
	"go-pizza-api/internal/pizzahut"
	"go-pizza-api/internal/ranking"
	"strconv"
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

	// Scoring
	Score float64

	// Discounts
	Discount            float64
	PriceBeforeDiscount float64
	PriceAfterDiscount  float64

	// Products
	InchesOfPizza float64
	Products      []ranking.Product
}

func GetDeals(postcode string) []AllDeals {
	deals := []AllDeals{}
	var wg sync.WaitGroup
	wg.Add(3)

	// Pizzahut
	go func() {
		defer wg.Done()
		pizzas, discounts, err := pizzahut.GetDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range pizzas {

			products := pizzahut.FormatProductData(item.DealContent)
			inchesOfPizza := ranking.CalculateTotalInches(products[0], ranking.PizzahutSizes)

			deals = append(deals, AllDeals{
				Restaurant: "Pizza Hut",
				DealName:   item.Title,
				DealDesc:   item.Desc,
				DealPrice:  item.Price,
				DealUrl:    "https://www.pizzahut.co.uk/order/deal/?id=" + item.ID,
				DealType:   item.Type,
				Products:   products,
				Score:      ranking.ScoreDeal(inchesOfPizza, item.Price),
			})
		}
		for _, disc := range discounts {

			score, costAfterDiscount := ranking.ScoreDiscount(disc.Reduction, disc.Price)

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
	go func() {
		defer wg.Done()
		savings, err := papajohns.GetDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range savings {
			if item.Displayed {

				var collection bool
				if item.ShippingMethod == 1 {
					collection = true
				}

				var products []ranking.Product
				var score float64
				var costAfterDiscount = -1.00

				if item.ReductionType == "%" {
					score, costAfterDiscount = ranking.ScoreDiscount(item.Reduction, item.MinimumSpend)
				} else {
					products = papajohns.FormatProductData(item.DealContent, item.Desc)

					pizzaType := ranking.LookupPizza(products)

					inchesOfPizza := ranking.CalculateTotalInches(pizzaType, ranking.PapajohsSizes)
					score = ranking.ScoreDeal(inchesOfPizza, item.Price)
				}

				deals = append(deals, AllDeals{
					Restaurant:          "Papa Johns",
					DealName:            item.DisplayName,
					DealUrl:             "https://www.papajohns.co.uk/deals",
					DealDesc:            item.Desc,
					DealPrice:           item.Price,
					DealType:            item.ReductionType,
					StudentDeal:         item.StudentDeal,
					CollectionOnly:      collection,
					Discount:            item.Reduction,
					PriceAfterDiscount:  float64(costAfterDiscount),
					PriceBeforeDiscount: item.MinimumSpend, //field only present on % deals
					Products:            products,
					Score:               score,
				})
			}
		}
	}()

	// Dominos
	go func() {
		defer wg.Done()
		domnios, vouchers, err := dominos.GetAllSavings(postcode)
		if err != nil {
			return
		}

		for _, item := range domnios[0].StoreDeals {

			products := dominos.FormatProductData(item.Deal[0].DealContent, item.Deal[0].Desc)
			pizzaType := ranking.LookupPizza(products)
			inchesOfPizza := ranking.CalculateTotalInches(pizzaType, ranking.PapajohsSizes)

			deals = append(deals, AllDeals{
				Restaurant:     "Domino's",
				DealName:       item.Name,
				DealUrl:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
				DealPrice:      item.Deal[0].Price,
				DealDesc:       item.Deal[0].Desc,
				DealType:       "",
				CollectionOnly: false,
				StudentDeal:    false, //TODO
				Products:       dominos.FormatProductData(item.Deal[0].DealContent, item.Deal[0].Desc),
				Score:          ranking.ScoreDeal(inchesOfPizza, item.Deal[0].Price),
			})
		}
		for _, item := range vouchers {
			reduction, err := dominos.GetReduction(item.Desc)
			if err != nil {
				return
			}

			score, costAfterDiscount := ranking.ScoreDiscount(reduction, item.MinSpend.Amount)

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
