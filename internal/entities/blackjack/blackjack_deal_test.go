package blackjack

import (
	"errors"
	"fmt"
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
		name               string
		dealerHand         *DealerHand
		expectedDealerHand *DealerHand
		expectedError      error
	}

	testCases := []testCase{
		{
			name:          "Dealer hand is nil",
			dealerHand:    nil,
			expectedError: errors.New("dealer hand is nil"),
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			actualErr := testBlackjack.DrawRemainingDealerCards(tc.dealerHand)

			if tc.expectedError != nil {
				require.EqualError(t, actualErr, tc.expectedError.Error())
				return
			}

			require.NoError(t, actualErr)
			require.Equal(t, tc.expectedDealerHand, tc.dealerHand)
		})
	}
}
