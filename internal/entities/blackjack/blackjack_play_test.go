package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/entities/player"
	"scrub/internal/testutils"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBlackjack_Play(t *testing.T) {
	type testCase struct {
		name            string
		inputPlayers    []BlackjackPlayer
		inputDealerHand DealerHand
		expectedErr     error
		expectedSplit   bool
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
			expectedSplit: false,
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
			expectedSplit: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, i, tc.name), func(t *testing.T) {
			testBlackjack := NewBlackjack(10)

			err = testBlackjack.Play(testLogger, tc.inputPlayers, tc.inputDealerHand, Strategy)

			if tc.expectedErr != nil {
				require.ErrorContains(t, tc.expectedErr, err.Error(), "error value")
				return
			}

			require.NoError(t, err, "unexpected error")

			if tc.expectedSplit {
				require.Equal(t, uint64(1), testBlackjack.SplitCount, "split count")
			}
		})
	}
}
