package blackjack

import (
	"errors"
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/testutils"
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

func TestBlackjack_DrawRemainingDealerCards(t *testing.T) {
	type testCase struct {
		name                       string
		dealerHand                 *DealerHand
		expectedError              error
		expectedNumberOfDrawnCards int
	}

	testCases := []testCase{
		{
			name:          "Dealer hand is nil",
			dealerHand:    nil,
			expectedError: errors.New("dealer hand is nil"),
		},
		{
			name: "Dealer hand has no value NoCards",
			dealerHand: &DealerHand{
				Hand: Hand{
					cards: []deck.Card{},
				},
			},
			expectedError: errors.New("dealer hand has no value yet"),
		},
		{
			name:          "Dealer hand has no value Nil",
			dealerHand:    &DealerHand{},
			expectedError: errors.New("dealer hand has no value yet"),
		},
		{
			name: "Dealer should draw 1 card",
			dealerHand: &DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.QueenOfHearts,
						deck.SixOfClubs,
					},
				},
			},
			expectedNumberOfDrawnCards: 1,
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			var beforeActiveCards, beforeDealerCards int
			
			if tc.dealerHand != nil {
				beforeActiveCards = len(testBlackjack.deck.ActiveCards)
				beforeDealerCards = len(tc.dealerHand.cards)
			}

			actualErr := testBlackjack.DrawRemainingDealerCards(tc.dealerHand)

			if tc.expectedError != nil {
				require.EqualError(t, actualErr, tc.expectedError.Error())
				return
			}

			require.NoError(t, actualErr)
			require.Len(t, tc.dealerHand.cards, beforeDealerCards+tc.expectedNumberOfDrawnCards)
			require.Len(t, testBlackjack.deck.ActiveCards, beforeActiveCards-tc.expectedNumberOfDrawnCards)
		})
	}
}
