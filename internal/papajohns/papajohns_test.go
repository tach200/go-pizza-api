package papajohns

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStoreInfo(t *testing.T) {
	// These may change, so even if the tests fails it may still work
	// check the test results.
	want := StoreInfo{
		// Version: 20221230141034,s
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

func TestFormatProductData(t *testing.T) {
	type args struct {
		dealContent []DealContent
	}
	tests := []struct {
		name string
		args args
		want []Product
	}{
		{
			name: "test that this works or something lol",
			args: args{
				dealContent: []DealContent{
					{Product: "Pizza"},
					{Product: "Pizza"},
					{Product: "Side"},
					{Product: "Side"},
					{Product: "Large Drink"},
				},
			},
			want: []Product{
				{
					ProductType:  "Pizza",
					ProductCount: 2,
				},
				{
					ProductType:  "Side",
					ProductCount: 2,
				},
				{
					ProductType:  "Large Drink",
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
