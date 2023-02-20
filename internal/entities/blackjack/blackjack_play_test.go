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
		name                  string
		inputPlayers          []BlackjackPlayer
		inputDealerHand       DealerHand
		expectedErr           error
		expectedNumberOfHands int
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
			inputPlayers: []BlackjackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								deck.TenOfClubs,
								deck.TenOfDiamonds,
							},
							BetAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.TenOfClubs,
						deck.TwoOfDiamonds,
						deck.TenOfClubs,
					},
				},
			},
			expectedNumberOfHands: 1,
		},
		{
			name: "OK Split",
			inputPlayers: []BlackjackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								deck.NineOfClubs,
								deck.NineOfDiamonds,
							},
							BetAmount: 10,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.AceOfClubs,
						deck.SixOfClubs,
					},
				},
			},
			expectedNumberOfHands: 2,
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
			require.Len(t, tc.inputPlayers[0].Hands, tc.expectedNumberOfHands, "number of hands")
		})
	}
}
