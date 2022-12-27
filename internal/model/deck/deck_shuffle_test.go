package deck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeck_Shuffle(t *testing.T) {
	testDeck := NewDeck()
	testDeck.TakeCardByIndex(0)
	testDeck.Shuffle()
	require.Len(t, testDeck.ActiveCards, 52, "burnt cards must be restored")
}
