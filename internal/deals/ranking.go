package deals

const (
	pricePattern      = "[+-]?([0-9]*[.])?[0-9]+"
	keywordPattern    = "(?m)(?i)\\bpersonal|\\bsmall|\\bmedium|\\blarge|\\bxxl|\\bside|\\badd|\\btoppings|([0-9] * [0-9]*)"
	percentagePattern = "(\\d+(\\.\\d+)?%)"
)

var (
	// sizes of pizzas in inches.
	dominosSizes = map[string]float64{
		"personal pizza": 7.0,
		"Personal":       7.0,
		"small pizza":    9.5,
		"Small":          9.5,
		"medium pizza":   11.5,
		"Medium":         11.5,
		"large pizza":    13.5,
		"Large":          13.5,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
	papajohsSizes = map[string]float64{
		"small pizza":  8,
		"Small":        8,
		"medium pizza": 11.5,
		"Medium":       11.5,
		"large pizza":  13.5,
		"Large":        13.5,
		"XXL":          15.5,
		"xxl pizza":    15.5,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
	pizzahutSizes = map[string]float64{
		"small pizza":  9,
		"Small":        9,
		"medium pizza": 11,
		"Medium":       11,
		"large pizza":  14,
		"Large":        14,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
)

type Product struct {
	ProductType  string
	ProductCount int
}

func scoreDiscount(reduction, cost float64) (float64, float64) {
	costAfterDiscount := costAfterDiscount(reduction, cost)
	score := reduction - cost

	return score, costAfterDiscount
}

func costAfterDiscount(reduction, cost float64) float64 {
	return cost - (cost * reduction / 100)
}

func calculateTotalInches(poduct Product, pizzaSizes map[string]float64) float64 {
	size, ok := pizzaSizes[poduct.ProductType]
	if !ok {
		return -1
	}

	return size * float64(poduct.ProductCount)
}

func scoreDeal(inchesOfPizza, price float64) float64 {
	return (inchesOfPizza / price)
}
