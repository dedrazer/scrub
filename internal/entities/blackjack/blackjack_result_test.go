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
		inputPlayers          []BlackJackPlayer
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
			inputPlayers: []BlackJackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
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
							betAmount: 50,
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
			expectedResultCredits: []uint64{1050},
		},
		{
			name: "OK Player Bust",
			inputPlayers: []BlackJackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
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
							betAmount: 50,
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
			expectedResultCredits: []uint64{950},
		},
		{
			name: "OK Player & Dealer Bust",
			inputPlayers: []BlackJackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
						{
							cards: []deck.Card{
								{
									Value:  7,
									Suit:   "Spades",
									Symbol: "7",
								},
								{
									Value:  4,
									Suit:   "Hearts",
									Symbol: "4",
								},
								{
									Value:  3,
									Suit:   "Spades",
									Symbol: "3",
								},
								{
									Value:  8,
									Suit:   "Clubs",
									Symbol: "8",
								},
							},
							betAmount: 50,
						},
					},
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{cards: []deck.Card{
					{
						Value:  6,
						Suit:   "Spades",
						Symbol: "6",
					},
					{
						Value:  8,
						Suit:   "Diamonds",
						Symbol: "8",
					},
					{
						Value:  8,
						Suit:   "Diamonds",
						Symbol: "8",
					},
				},
				},
			},
			expectedResultCredits: []uint64{950},
		},
		{
			name: "OK Push",
			inputPlayers: []BlackJackPlayer{
				{
					Player: testPlayer,
					Hands: []Hand{
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
							betAmount: 50,
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
