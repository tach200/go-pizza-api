package deals

import (
	"reflect"
	"testing"
)

func Test_getDealKeywords(t *testing.T) {
	type args struct {
		dealDesc string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "extract correct keywords from the deal",
			args: args{
				dealDesc: "X1 Large Stuffed Crust Pizzas And A Classic Side",
			},
			want:    []string{"Large", "Side"},
			wantErr: false,
		},
		{
			name: "extract correct keywords from the deal",
			args: args{
				dealDesc: "2 sumptuous sides for £5.49 - Choose between garlic pizza bread and potato wedges.",
			},
			want:    []string{"2", "Sides"},
			wantErr: false,
		},
		{
			name: "extract correct keywords from the deal",
			args: args{
				dealDesc: "Any Large Pizza, 7-piece Chicken Side, Garlic Pizza Bread and 1.25l drink for £27.99. Choose a pizza from the menu or create your own up to 4 toppings. Premium crusts and additional toppings will be charged as extra.",
			},
			want:    []string{"Large", "Side"},
			wantErr: false,
		},
		{
			name: "extract correct keywords from the deal",
			args: args{
				dealDesc: "7 wings or 7 chunks, wedges, 500ml drink",
			},
			want:    []string{"7 ", "7 "},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDealKeywords(tt.args.dealDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDealKeywords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDealKeywords() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_rankAllPapa(t *testing.T) {
// 	deals, err := papajohns.GetDeals("ME46EA")
// 	assert.Nil(t, err)

// 	for _, deal := range deals {
// 		score := rankScore(deal.Desc, deal.Price, papajohsSizes)

// 		assert.NotEqual(t, float64(-1), score)
// 		assert.NotEqual(t, float64(0), score)
// 	}
// }

// func Test_rank(t *testing.T) {
// 	desc := "Any medium pizza, 1 classic side and any regular drink for only £17.99"
// 	price := 4.99

// 	score := rankScore(desc, price, papajohsSizes)

// 	assert.NotEqual(t, float64(-1), score)
// 	assert.NotEqual(t, float64(0), score)
// }

// func Test_rankPizzaHut(t *testing.T) {
// 	deals, err := pizzahut.GetDeals("ME46EA")
// 	assert.Nil(t, err)

// 	for _, deal := range deals {
// 		score := rankScore(deal.Desc, deal.Price, pizzahutSizes)

// 		assert.NotEqual(t, float64(-1), score)
// 		assert.NotEqual(t, float64(0), score)
// 	}
// }

func Test_calculateScoreArr(t *testing.T) {
	type args struct {
		scoreArr []keywordType
		dealCost float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// {
		// 	name: "calculate the score of the array",
		// 	args: args{
		// 		scoreArr: []keywordType{
		// 			{
		// 				value:      2,
		// 				multiplier: true,
		// 			},
		// 			{
		// 				value:      0,
		// 				multiplier: false,
		// 			},
		// 			{
		// 				value:      9,
		// 				multiplier: true,
		// 			},
		// 		},
		// 	},
		// },
		{
			name: "calculate the score of the array",
			args: args{
				scoreArr: []keywordType{
					{
						value:      7,
						multiplier: true,
					},
					{
						value:      7,
						multiplier: true,
					},
				},
				dealCost: 10.99,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateScoreArr(tt.args.scoreArr, tt.args.dealCost)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateScoreArr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateScoreArr() = %v, want %v", got, tt.want)
			}
		})
	}
}
