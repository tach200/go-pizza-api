package deals

import (
	"errors"
	"regexp"
	"strconv"
)

const (
	pricePattern   = "[+-]?([0-9]*[.])?[0-9]+"
	keywordPattern = "(?i)\\bpersonal|\\bsmall|\\bmedium|\\blarge|\\bxxl|([0-9])"
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
	}
	pizzahutSizes = map[string]float64{
		"small":  9,
		"Small":  9,
		"medium": 11,
		"Medium": 11,
		"large":  14,
		"Large":  14,
	}
)

// rankScore attempts to generate a value which reflects how good the deal is
// higher is better, but this is subject to change in the future.
func rankScore(dealTitle, dealDesc string, pizzaSizes map[string]float64) float64 {
	keywords, err := getDealKeywords(dealDesc)
	if err != nil {
		return -1
	}

	scoreArr, err := convertToScoreArr(keywords, pizzaSizes)
	if err != nil {
		return -1
	}

	dealCost, err := getDealCost(dealTitle, dealDesc)
	if err != nil {
		return -1
	}

	score, err := calculateScoreArr(scoreArr, dealCost)
	if err != nil {
		return -1
	}

	return score
}

// getDealCost will attempt to return the cost of the deal
// the information that is returned is not guranteed to be the cost, but is most likely
// a regular expression is used to find numbers
// the cost doesn't always include a '£' symbol, which makes this process a bit more complex.
func getDealCost(dealTitle, dealDesc string) (float64, error) {
	reg := regexp.MustCompile(pricePattern)
	titleCosts := reg.FindAllString(dealTitle, -1)
	descCosts := reg.FindAllString(dealDesc, -1)

	// Combine the two lists before converting to floats
	floats, err := costStrsToFloats(append(titleCosts, descCosts...))

	if err != nil {
		return 0, err
	}

	// The largest float number in the list is most likely to be the cost
	return findMaxFloat(floats), nil
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
func convertToScoreArr(keywords []string, pizzaSizes map[string]float64) ([]float64, error) {
	len := len(keywords)
	var floats []float64
	switch {
	case len == 1:
		floats = append(floats, pizzaSizes[keywords[0]])
	case len >= 2:

		amount, err := strconv.ParseFloat(keywords[0], 32)
		if err != nil {
			return nil, err
		}

		floats = append(floats, amount)
		floats = append(floats, pizzaSizes[keywords[1]])
	default:
		return nil, errors.New("error: unimplemented or an error")
	}

	return floats, nil
}

func calculateScoreArr(scoreArr []float64, dealCost float64) (float64, error) {
	len := len(scoreArr)

	switch len {
	case 1:
		return (scoreArr[0] / dealCost), nil
	case 2:
		return ((scoreArr[0] * scoreArr[1]) / dealCost), nil
	default:
		return -1, errors.New("error: unimplemented or an error")
	}
}

func costStrsToFloats(strFloats []string) ([]float64, error) {
	var floats []float64

	for _, strFl := range strFloats {
		float, err := strconv.ParseFloat(strFl, 64)
		if err != nil {
			return nil, err
		}
		floats = append(floats, float)
	}

	return floats, nil
}

func findMaxFloat(floats []float64) (max float64) {
	max = 0.00
	for _, float := range floats {
		if float > max {
			max = float
		}
	}
	return max
}
