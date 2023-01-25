package blackjack

import (
	"scrub/internal/entities/deck"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBlackjack_Result(t *testing.T) {
	type testCase struct {
		name             string
		inputPlayerHands []Hand
		inputDealerHand  DealerHand
		expectedResults  []Result
		expectedErr      error
	}

	testCases := []testCase{
		{
			name: "OK Dealer Bust",
			inputPlayerHands: []Hand{
				{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "King",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Queen",
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "Queen",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Jack",
						},
						{
							Value:  2,
							Suit:   "Hearts",
							Symbol: "2",
						},
					},
				},
			},
			expectedResults: []Result{
				{
					Status: &win,
				},
			},
		},
		{
			name: "OK Player Bust",
			inputPlayerHands: []Hand{
				{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "King",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Queen",
						},
						{
							Value:  2,
							Suit:   "Hearts",
							Symbol: "2",
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "Queen",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Jack",
						},
					},
				},
			},
			expectedResults: []Result{
				{
					Status: &loss,
				},
			},
		},
		{
			name: "OK Push",
			inputPlayerHands: []Hand{
				{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "King",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Queen",
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "Queen",
						},
						{
							Value:  10,
							Suit:   "Diamonds",
							Symbol: "Jack",
						},
					},
				},
			},
			expectedResults: []Result{
				{
					Status: nil,
				},
			},
		},
	}

	var numberOfDecks uint = 4
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "zap.NewDevelopment() setup error")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bj := NewBlackjack(numberOfDecks)
			var results []Result
			results, err = bj.Results(logger, tc.inputPlayerHands, tc.inputDealerHand)

			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error(), "expected error message")
				return
			}

			require.NoError(t, err, "unexpected error")
			require.EqualValues(t, tc.expectedResults, results, "expected results")
		})
	}
}
