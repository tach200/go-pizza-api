package deals

import (
	"go-pizza-api/internal/dominos"
	"go-pizza-api/internal/papajohns"
	"go-pizza-api/internal/pizzahut"
	"strconv"
	"sync"
)

type AllDeals struct {
	Restaurant string
	DealName   string
	DealDesc   string
	Url        string
	// Rank is currently not supported.
	Rank float64
}

func GetDeals(postcode string) []AllDeals {

	// span := tracer.StartSpan("/GET Deals")
	// defer span.Finish()

	//Create list of structs to store clean data.
	deals := []AllDeals{}
	var wg sync.WaitGroup
	// 3 requests need to be made
	wg.Add(3)

	go func() {
		// pizzaHutSpan := tracer.StartSpan("/GET PizzaHut", opentracing.ChildOf(span.Context()))
		// defer pizzaHutSpan.Finish()

		defer wg.Done()
		pizzahut, err := pizzahut.GetPizzahutDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range pizzahut {
			deals = append(deals, AllDeals{
				Restaurant: "Pizza Hut",
				DealName:   item.Title,
				DealDesc:   item.Desc,
				Url:        "https://www.pizzahut.co.uk/order/deal/?id=" + item.Id,
				Rank:       rankScore(item.Title, item.Desc, pizzahutSizes),
			})
		}
	}()

	go func() {
		// papaJohnSpan := tracer.StartSpan("/GET PapaJohn", opentracing.ChildOf(span.Context()))
		// defer papaJohnSpan.Finish()

		defer wg.Done()
		papajohns, err := papajohns.GetPapajohnsDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range papajohns {
			deals = append(deals, AllDeals{
				Restaurant: "Papa Johns",
				DealName:   item.Name,
				Url:        "https://www.papajohns.co.uk/" + item.Url,
				DealDesc:   item.Desc,
				Rank:       rankScore(item.Name, item.Desc, papajohsSizes),
			})
		}
	}()

	go func() {
		// dominosSpan := tracer.StartSpan("/GET Dominos", opentracing.ChildOf(span.Context()))
		// defer dominosSpan.Finish()

		defer wg.Done()
		domDeals, vouchers, err := dominos.GetDominosDealsVouchers(postcode)
		if err != nil {
			return
		}
		for _, item := range domDeals[0].StoreDeals {
			deals = append(deals, AllDeals{
				Restaurant: "Domino's",
				DealName:   item.Name,
				Url:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
				DealDesc:   item.Deal[0].Desc,
				Rank:       rankScore(item.Name, item.Deal[0].Desc, dominosSizes),
			})
		}
		for _, item := range vouchers {
			deals = append(deals, AllDeals{
				Restaurant: "Domino's",
				DealName:   "Savings Voucher",
				Url:        "https://www.dominos.co.uk/deals",
				DealDesc:   item.Desc,
				Rank:       rankScore("Savings Voucher", item.Desc, dominosSizes),
			})
		}
	}()

	wg.Wait()
	return deals
}
