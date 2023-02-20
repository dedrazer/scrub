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

func TestBlackjack_Result(t *testing.T) {
	type testCase struct {
		name                  string
		inputPlayers          []BlackjackPlayer
		inputDealerHand       DealerHand
		expectedResultCredits []uint64
		expectedErr           error
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
								deck.KingOfClubs,
								deck.QueenOfDiamonds,
							},
							BetAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.QueenOfClubs,
						deck.JackOfDiamonds,
						deck.TwoOfHearts,
					},
				},
			},
			expectedResultCredits: []uint64{1050},
		},
		{
			name: "OK Player Bust",
			inputPlayers: []BlackjackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								deck.KingOfClubs,
								deck.QueenOfDiamonds,
								deck.TwoOfHearts,
							},
							BetAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.QueenOfClubs,
						deck.JackOfDiamonds,
					},
				},
			},
			expectedResultCredits: []uint64{950},
		},
		{
			name: "OK Player & Dealer Bust",
			inputPlayers: []BlackjackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								deck.SevenOfSpades,
								deck.FourOfHearts,
								deck.ThreeOfSpades,
								deck.EightOfClubs,
							},
							BetAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{cards: []deck.Card{
					deck.SixOfSpades,
					deck.EightOfDiamonds,
					deck.EightOfDiamonds,
				},
				},
			},
			expectedResultCredits: []uint64{950},
		},
		{
			name: "OK Push",
			inputPlayers: []BlackjackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								deck.KingOfClubs,
								deck.QueenOfDiamonds,
							},
							BetAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.QueenOfClubs,
						deck.JackOfDiamonds,
					},
				},
			},
			expectedResultCredits: []uint64{1000},
		},
	}

	var numberOfDecks uint = 4
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "zap.NewDevelopment() setup error")
	for i, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, i, tc.name), func(t *testing.T) {
			bj := NewBlackjack(numberOfDecks)
			err = bj.Results(logger, tc.inputPlayers, tc.inputDealerHand)

			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error(), "expected error message")
				return
			}

			require.NoError(t, err, "unexpected error")
			for j, inputPlayer := range tc.inputPlayers {
				require.Equal(t, tc.expectedResultCredits[j], inputPlayer.Player.Credits, "player credits. expected: %d, actual: %d.", tc.expectedResultCredits[j], inputPlayer.Player.Credits)
			}
		})
	}
}
