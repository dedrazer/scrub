package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/entities/player"
	internalTesting "scrub/internal/testing"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBlackjack_Play(t *testing.T) {
	type testCase struct {
		name            string
		inputPlayers    []BlackJackPlayer
		inputDealerHand DealerHand
		expectedErr     error
	}

	testLogger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to init zap logger")

	testPlayer := player.Player{
		Name:    "Player 1",
		Credits: 1000,
	}

	testCases := []testCase{
		{
			name: "OK Dealer Bust",
			inputPlayers: []BlackJackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								{
									Value:  10,
									Suit:   "Clubs",
									Symbol: "10",
								},
								{
									Value:  10,
									Suit:   "Diamonds",
									Symbol: "10",
								},
							},
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
							Symbol: "10",
						},
						{
							Value:  2,
							Suit:   "Diamonds",
							Symbol: "2",
						},
						{
							Value:  10,
							Suit:   "Clubs",
							Symbol: "10",
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, i, tc.name), func(t *testing.T) {
			testBlackjack := NewBlackjack(10)

			err = testBlackjack.Play(testLogger, tc.inputPlayers, tc.inputDealerHand, Strategy1)

			if tc.expectedErr != nil {
				require.ErrorContains(t, tc.expectedErr, err.Error(), "error value")
				return
			}

			require.NoError(t, err, "unexpected error")
		})
	}
}
