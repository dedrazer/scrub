package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck_Swap(t *testing.T) {
	testDeck := NewDeck()
	testDeck.Swap(0, 1)
	a, err := testDeck.GetCardByIndex(0)
	require.NoError(t, err, "must not return error")
	b, err := testDeck.GetCardByIndex(1)
	require.NoError(t, err, "must not return error")

	assert.Equal(t, a.Print(), "2 of Clubs", "first card")
	assert.Equal(t, b.Print(), "Ace of Clubs", "second card")
}
