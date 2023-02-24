package papajohns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStoreInfo(t *testing.T) {
	want := StoreInfo{
		Data: StoreID{
			ID: 440,
		},
	}

	got, err := getStoreID("me46ea")
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGetDeals(t *testing.T) {

	want := Deals{
		Deals: []Deal{
			{},
		},
	}

	got, err := GetDeals("me46ea")
	assert.Nil(t, err)
	assert.Equal(t, want, got)

}

// func TestFormatProductData(t *testing.T) {
// 	type args struct {
// 		dealContent []DealContent
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []ranking.Product
// 	}{
// 		{
// 			name: "test that this works or something lol",
// 			args: args{
// 				dealContent: []DealContent{
// 					{Product: "Pizza"},
// 					{Product: "Pizza"},
// 					{Product: "Side"},
// 					{Product: "Side"},
// 					{Product: "Large Drink"},
// 				},
// 			},
// 			want: []ranking.Product{
// 				{
// 					ProductType:  "Pizza",
// 					ProductCount: 2,
// 				},
// 				{
// 					ProductType:  "Side",
// 					ProductCount: 2,
// 				},
// 				{
// 					ProductType:  "Large Drink",
// 					ProductCount: 1,
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := FormatProductData(tt.args.dealContent); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("FormatProductData() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
