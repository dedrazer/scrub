package deck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeck_Merge(t *testing.T) {
	testDeck1 := NewDeck()
	testDeck2 := NewDeck()
	res := Merge(&testDeck1, &testDeck2)
	require.Len(t, res.ActiveCards, 104, "2 decks must have 104 cards")
}
