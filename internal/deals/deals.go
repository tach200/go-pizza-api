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
	//Create list of structs to store clean data.
	deals := []AllDeals{}
	var wg sync.WaitGroup
	// 3 requests need to be made
	wg.Add(3)

	go func() {
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
				Rank:       0,
			})
		}
	}()

	go func() {
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
				Rank:       0,
			})
		}
	}()

	go func() {
		defer wg.Done()
		dominos, err := dominos.GetDominosDeals(postcode)
		if err != nil {
			return
		}
		for _, item := range dominos[0].StoreDeals {
			deals = append(deals, AllDeals{
				Restaurant: "Domino's",
				DealName:   item.Name,
				Url:        "https://www.dominos.co.uk/deals/deal/" + strconv.Itoa(item.Deal[0].Id),
				DealDesc:   item.Deal[0].Desc,
				Rank:       0,
			})
		}
	}()

	wg.Wait()
	return deals
}
