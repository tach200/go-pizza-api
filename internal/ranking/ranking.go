package ranking

import (
	"math"
	"regexp"
	"strings"
)

var (
	// Sizes of pizzas in inches.
	DominosSizes = map[string]float64{
		"personal pizza": 7.0,
		"small pizza":    9.5,
		"medium pizza":   11.5,
		"large pizza":    13.5,
	}
	PapajohsSizes = map[string]float64{
		"small pizza":  8,
		"medium pizza": 11.5,
		"large pizza":  13.5,
		"xxl pizza":    15.5,
	}
	PizzahutSizes = map[string]float64{
		"small pizza":  9,
		"medium pizza": 11,
		"large pizza":  14,
	}
)

const pizzaSizeRegx = "(?m)(?i)\\bpersonal|\\bsmall|\\bmedium|\\blarge|\\bxxl"

// Product describes product data and count
type Product struct {
	ProductType  string
	ProductCount int
}

// ScoreDiscount returns a score for percentage based deals.
func ScoreDiscount(reduction, cost float64) (float64, float64) {
	costAfterDiscount := costAfterDiscount(reduction, cost)
	score := reduction - cost

	return score, costAfterDiscount
}

func costAfterDiscount(reduction, cost float64) float64 {
	return cost - (cost * reduction / 100)
}

// CalculateTotalInches calculates how many inches of pizza are included in a deal.
func CalculateTotalInches(poduct Product, pizzaSizes map[string]float64) float64 {
	size, ok := pizzaSizes[poduct.ProductType]
	if !ok {
		return -1
	}

	return size * float64(poduct.ProductCount)
}

// ScoreDeal calculates the cost per inch of pizza, which is used
// as a scoring function.
func ScoreDeal(inchesOfPizza, price float64) float64 {

	result := roundFloat((inchesOfPizza / price), 5)

	if math.IsInf(result, 1) {
		return -1
	}

	if math.IsInf(result, -1) {
		return -1
	}

	if result < 0 {
		return -1
	}

	return result
}

// GetPizzaSize finds out what size pizza is included in this deal
// this is useful when api can return ambiguous data.
func GetPizzaSize(desc string) string {
	regx := regexp.MustCompile(pizzaSizeRegx)

	size := regx.FindString(desc)
	if size == "" {
		return "unknown"
	}

	return size
}

// LookupPizza will return the pizza product in the list.
func LookupPizza(products []Product) Product {
	for _, prod := range products {
		if strings.Contains(prod.ProductType, "pizza") {
			return prod
		}
	}
	return Product{}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
