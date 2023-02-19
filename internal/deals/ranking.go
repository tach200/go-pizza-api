package deals

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	pricePattern      = "[+-]?([0-9]*[.])?[0-9]+"
	keywordPattern    = "(?m)(?i)\\bpersonal|\\bsmall|\\bmedium|\\blarge|\\bxxl|\\bside|\\badd|\\btoppings|([0-9] * [0-9]*)"
	percentagePattern = "(\\d+(\\.\\d+)?%)"
)

var (
	// sizes of pizzas in inches.
	dominosSizes = map[string]float64{
		"personal": 7.0,
		"Personal": 7.0,
		"small":    9.5,
		"Small":    9.5,
		"medium":   11.5,
		"Medium":   11.5,
		"large":    13.5,
		"Large":    13.5,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
	papajohsSizes = map[string]float64{
		"small":  8,
		"Small":  8,
		"medium": 11.5,
		"Medium": 11.5,
		"large":  13.5,
		"Large":  13.5,
		"XXL":    15.5,
		"xxl":    15.5,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
	pizzahutSizes = map[string]float64{
		"small":  9,
		"Small":  9,
		"medium": 11,
		"Medium": 11,
		"large":  14,
		"Large":  14,
		// Used to filter out junk
		"side":     0,
		"Side":     0,
		"toppings": 0,
		"Toppings": 0,
	}
)

type keywordType struct {
	value      float64
	multiplier bool
}

// rankScore attempts to generate a value which reflects how good the deal iszoo
// higher is better, but this is subject to change in the future.
func rankDealScore(dealDesc string, cost float64, pizzaSizes map[string]float64) float64 {
	keywords, err := getDealKeywords(dealDesc)
	if err != nil {
		// fmt.Printf("error: error finding deal keywords %s", err)
		return -1
	}

	scoreArr, err := convertToScoreArr(keywords, pizzaSizes)
	// fmt.Printf("error: error converting to score array %s", err)
	if err != nil {
		return -1
	}

	dealCost := cost
	if cost == 0 {
		// fmt.Println("error: cost is 0")
		return -1 //for now
	}

	score, err := calculateScoreArr(scoreArr, dealCost)
	if err != nil {
		// fmt.Printf("error: couldnt calculate a score %s", err)
		return -1
	}

	return score
}

// getDealKeywords will extracts keywords that will be used to help calculate the final score.
func getDealKeywords(dealDesc string) ([]string, error) {
	reg := regexp.MustCompile(keywordPattern)

	keywords := reg.FindAllString(dealDesc, -1)

	if len(keywords) < 1 {
		return nil, errors.New("error: no keywords extracted from text")
	}

	return keywords, nil
}

// convertToScoreArr converts a text string such as 'large' into the corresponding pizza size in inches.
func convertToScoreArr(keywords []string, pizzaSizes map[string]float64) ([]keywordType, error) {

	keywordValues := make([]keywordType, 0)

	for _, keyword := range keywords {
		// TODO: better handle this keyword for deals
		// don't want any keywords after add
		if keyword == "add" {
			break
		}

		val, ok := pizzaSizes[keyword]
		if ok {

			keywordType := keywordType{
				value:      val,
				multiplier: false,
			}

			keywordValues = append(keywordValues, keywordType)
			continue
		}

		formatStr := strings.TrimSpace(keyword)
		multiplier, err := strconv.ParseFloat(formatStr, 64)
		if err != nil {
			// fmt.Printf("error: couldn't parse float %s", err)
			continue
		}

		keywordType := keywordType{
			value:      multiplier,
			multiplier: true,
		}

		keywordValues = append(keywordValues, keywordType)
	}

	return keywordValues, nil
}

// calculateScoreArr calculates the final score.
func calculateScoreArr(scoreArr []keywordType, dealCost float64) (float64, error) {
	if dealCost == 0 {
		return -1, errors.New("error: couldn't rank deal")
	}

	multiplyNextVal := false
	valueToMultiply := 0.00
	inchesOfPizza := 0.00

	for _, val := range scoreArr {
		if multiplyNextVal && !val.multiplier {
			inchesOfPizza += (valueToMultiply * val.value)
			valueToMultiply = 0
			continue
		}

		if val.multiplier {
			multiplyNextVal = true
			valueToMultiply = val.value
			continue
		}

		inchesOfPizza += val.value
	}

	return (inchesOfPizza / dealCost), nil
}

func scoreDiscount(reduction, cost float64) (float64, float64) {
	costAfterDiscount := costAfterDiscount(reduction, cost)
	score := reduction - cost

	return score, costAfterDiscount
}

func costAfterDiscount(reduction, cost float64) float64 {
	return cost - (cost * reduction / 100)
}
