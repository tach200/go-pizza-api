package deals

import (
	"go-pizza-api/internal/dominos"
	"go-pizza-api/internal/papajohns"
	"go-pizza-api/internal/pizzahut"
	"go-pizza-api/internal/ranking"
	"strconv"
	"sync"
)

type Deals struct {
	Deals    []AllDeals
	Postcode string
}
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

func (d *Deals) GetPizzahutDeals() {
	pizzas, discounts, err := pizzahut.GetDeals(d.Postcode)
	if err != nil {
		return
	}
	for _, item := range pizzas {

		products := pizzahut.FormatProductData(item.DealContent)
		inchesOfPizza := ranking.CalculateTotalInches(products[0], ranking.PizzahutSizes)

		d.Deals = append(d.Deals, AllDeals{
			Restaurant:    "Pizza Hut",
			DealName:      item.Title,
			DealDesc:      item.Desc,
			DealPrice:     item.Price,
			DealUrl:       "https://www.pizzahut.co.uk/order/deal/?id=" + item.ID,
			DealType:      "items",
			Products:      products,
			InchesOfPizza: inchesOfPizza,
			StudentDeal:   pizzahut.IsStudentDeal(item.ID),
			Score:         ranking.ScoreDeal(inchesOfPizza, item.Price),
		})
	}
	for _, disc := range discounts {
		score, costAfterDiscount := ranking.ScoreDiscount(disc.Reduction, disc.Price)

		d.Deals = append(d.Deals, AllDeals{
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
}

func (d *Deals) GetPapajohnsDeals() {
	savings, err := papajohns.GetDeals(d.Postcode)
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

			d.Deals = append(d.Deals, AllDeals{
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
}

func (d *Deals) GetDominosDeals() {
	doms, vouchers, err := dominos.GetAllSavings(d.Postcode)
	if err != nil {
		return
	}

	for _, storeDeal := range doms {
		for _, item := range storeDeal.StoreDeals {

			products := dominos.FormatProductData(item.Deal[0].DealContent, item.Deal[0].Desc)
			pizzaType := ranking.LookupPizza(products)
			inchesOfPizza := ranking.CalculateTotalInches(pizzaType, ranking.PapajohsSizes)

			d.Deals = append(d.Deals, AllDeals{
				Restaurant:     "Domino's",
				DealName:       item.Name,
				DealUrl:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
				DealPrice:      item.Deal[0].Price,
				DealDesc:       item.Deal[0].Desc,
				DealType:       "product",
				CollectionOnly: false,
				StudentDeal:    false, // TODO
				Products:       dominos.FormatProductData(item.Deal[0].DealContent, item.Deal[0].Desc),
				Score:          ranking.ScoreDeal(inchesOfPizza, item.Deal[0].Price),
			})
		}
	}

	for _, item := range vouchers {
		reduction, err := dominos.GetReduction(item.Desc)
		if err != nil {
			return
		}

		score, costAfterDiscount := ranking.ScoreDiscount(reduction, item.MinSpend.Amount)

		d.Deals = append(d.Deals, AllDeals{
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
}

func GetDeals(postcode string) []AllDeals {
	deals := Deals{
		Postcode: postcode,
		Deals:    []AllDeals{},
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		deals.GetDominosDeals()
	}()

	go func() {
		defer wg.Done()
		deals.GetPapajohnsDeals()
	}()

	go func() {
		defer wg.Done()
		deals.GetPizzahutDeals()

	}()

	wg.Wait()
	return deals.Deals
}
