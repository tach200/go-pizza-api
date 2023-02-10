package pizzahut

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeals(t *testing.T) {
	postcode := "ME46EA"

	deals, err := GetDeals(postcode)
	assert.Nil(t, err)
	assert.IsType(t, []MenuItem{}, deals, "")
}
