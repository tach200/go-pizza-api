package papajohns

import (
	"errors"
	"go-pizza-api/internal/request"

	"github.com/gocolly/colly"
)

// No API exposed so have to webscrape for this one.
func papajohnsStoreInfo(postcode string) string {
	// Create the endpoint
	endpoint := "https://www.papajohns.co.uk/store-locator.aspx?postcode=" + postcode

	// Scrape the store URL from the store-locator to ensure correct deals.
	c := colly.NewCollector()
	var dealsUrl string
	c.OnHTML("#ctl00_cphBody_hypStoreDetails", func(element *colly.HTMLElement) {
		storeURL := element.Attr("href")
		dealsUrl = storeURL
	})
	err := c.Visit(endpoint)
	if err != nil {
		return ""
	}

	// Data has some junk at the end so trim it.
	// fmt.Println(dealsUrl[:len(dealsUrl)-9])
	// log.Print(dealsUrl[:len(dealsUrl)-9])
	return dealsUrl[:len(dealsUrl)-9]
}

type papajohnsDeal struct {
	Name string
	Desc string
	Url  string
}

func GetPapajohnsDeals(postcode string) ([]papajohnsDeal, error) {
	deals := []papajohnsDeal{}

	// Construct endpoint
	endpoint := "https://www.papajohns.co.uk" + papajohnsStoreInfo(postcode) + "offers.aspx"
	c := colly.NewCollector()
	c.UserAgent = request.UserAgent
	// Some deals are in these offer body tags.
	c.OnHTML(".menuListCont", func(element *colly.HTMLElement) {
		// deal := element.Text
		deals = append(deals, papajohnsDeal{
			Name: element.ChildText(".w100"),
			Desc: element.ChildText("p"),
			//This element has two a attributes, the second one is the one that hold the url that is needed.
			Url: element.ChildAttrs("a", "href")[1],
		})
	})
	// Some other deals are in this text tag
	c.OnHTML(".moreOffersList", func(element *colly.HTMLElement) {
		deals = append(deals, papajohnsDeal{
			Name: element.ChildText("h3.offer--title"),
			Desc: element.ChildText("p"),
			Url:  element.ChildAttr("a", "href"),
		})
	})
	err := c.Visit(endpoint)
	if err != nil {
		return nil, err
	}

	if len(deals) == 0 {
		return deals, errors.New("no deals found")
	}

	return deals, nil
}
