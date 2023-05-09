package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/entities/player"
	"scrub/internal/testutils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlackjack_Play(t *testing.T) {
	type testCase struct {
		name            string
		inputPlayers    []BlackjackPlayer
		inputDealerHand DealerHand
		expectedErr     error
		expectedSplit   bool
	}

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
			resetBJ()

			err := testBlackjack.Play(tc.inputPlayers, tc.inputDealerHand)

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

func TestBlackjack_double(t *testing.T) {
	testPlayer := BlackjackPlayer{
		Hands: []Hand{
			{
				cards: []deck.Card{
					deck.QueenOfHearts,
				},
				isDoubled: false,
				BetAmount: 10,
			},
			{
				cards: []deck.Card{
					deck.ThreeOfClubs,
				},
				isDoubled: false,
				BetAmount: 50,
			},
		},
	}

	err := testBlackjack.double(&testPlayer, 0)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, uint64(20), testPlayer.Hands[0].BetAmount, "bet amount doubled")
	require.True(t, testPlayer.Hands[0].isDoubled, "is doubled")

	require.Equal(t, uint64(50), testPlayer.Hands[1].BetAmount, "bet amount unchanged")
	require.False(t, testPlayer.Hands[1].isDoubled, "is not doubled")
}

func TestBlackjack_double_Error(t *testing.T) {
	testPlayer := BlackjackPlayer{
		Hands: []Hand{
			{
				cards: []deck.Card{
					deck.QueenOfHearts,
				},
				isDoubled: false,
				BetAmount: 10,
			},
			{
				cards: []deck.Card{
					deck.ThreeOfClubs,
				},
				isDoubled: false,
				BetAmount: 50,
			},
		},
	}

	testBlackjack.deck.ActiveCards = []deck.Card{}

	err := testBlackjack.double(&testPlayer, 0)

	require.EqualError(t, err, "Failed to dealPlayerACard: Failed to DealCard: Failed to TakeCardByIndex: Index is out of range")
}
