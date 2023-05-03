package blackjack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlackjack_DealCard(t *testing.T) {
	expectedActiveCards := len(testBlackjack.deck.ActiveCards) - 1
	expectedBurntCards := len(testBlackjack.deck.BurntCards) + 1

	card, err := testBlackjack.DealCard()
	if err != nil {
		t.Fatalf("Failed to deal card: %s", err.Error())
	}

	require.NotNil(t, card)
	require.Len(t, testBlackjack.deck.ActiveCards, expectedActiveCards)
	require.Len(t, testBlackjack.deck.BurntCards, expectedBurntCards)
}
