package blackjack

import (
	"fmt"
	"scrub/internal/entities/blackjack/utils"
	"scrub/internal/entities/deck"
	"scrub/internal/entities/player"
	"scrub/internal/testutils"
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
		expectedResult        string
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
			expectedResult:        utils.Win,
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
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, i, tc.name), func(t *testing.T) {
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

			switch tc.expectedResult {
			case utils.Win:
				require.Equal(t, uint64(1), bj.PlayerWins)
				require.Equal(t, uint64(0), bj.PlayerLosses)
				require.Equal(t, uint64(0), bj.Pushes)
				require.Equal(t, uint64(0), bj.PlayerBlackjackCount)
				require.Equal(t, uint64(0), bj.PlayerBust)

				for k := range tc.inputPlayers {
					require.Equal(t, uint64(1), tc.inputPlayers[k].Wins)
					require.Equal(t, uint64(0), tc.inputPlayers[k].Losses)
					require.Equal(t, uint64(0), tc.inputPlayers[k].Draws)
				}
			default:
			}
		})
	}
}
