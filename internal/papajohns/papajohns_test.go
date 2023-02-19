package papajohns

import (
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
