package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck_TakeCardByIndex(t *testing.T) {
	testDeck := NewDeck()

	card, err := testDeck.TakeCardByIndex(0)
	require.NoError(t, err, "must not return error")
	assert.Equal(t, uint(11), card.Value, "should return card with value 2")
	assert.Len(t, testDeck.ActiveCards, 51, "51 cards should remain")
	assert.Len(t, testDeck.BurntCards, 1, "1 card should be burnt")
}
