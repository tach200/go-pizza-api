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
