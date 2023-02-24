package pizzahut

import (
	"go-pizza-api/internal/ranking"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeals(t *testing.T) {
	postcode := "ME46EA"

	pizzas, vouchers, err := GetDeals(postcode)
	assert.Nil(t, err)
	assert.IsType(t, []MenuItem{}, pizzas, "")
	assert.IsType(t, []MenuItem{}, vouchers, "")
}

func TestFormatProductData(t *testing.T) {
	type args struct {
		dealContent []DealContent
	}
	tests := []struct {
		name string
		args args
		want []ranking.Product
	}{
		{
			name: "test lol",
			args: args{
				dealContent: []DealContent{
					{
						Count:     2,
						PizzaSize: "Medium",
						Product:   "pizza",
					},
					{
						Count:     1,
						PizzaSize: "Medium",
						Product:   "pizza",
					},
					{
						Count:     2,
						PizzaSize: "",
						Product:   "side",
					},
					{
						Count:     1,
						PizzaSize: "",
						Product:   "side",
					},
					{
						Count:     1,
						PizzaSize: "",
						Product:   "drink",
					},
				},
			},
			want: []ranking.Product{
				{
					ProductType:  "medium pizza",
					ProductCount: 3,
				},
				{
					ProductType:  "side",
					ProductCount: 3,
				},
				{
					ProductType:  "drink",
					ProductCount: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatProductData(tt.args.dealContent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatProductData() = %v, want %v", got, tt.want)
			}
		})
	}
}
