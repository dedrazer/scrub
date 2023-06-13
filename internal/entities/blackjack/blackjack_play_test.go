package blackjack

import (
	"errors"
	"fmt"
	bjutils "scrub/internal/entities/blackjack/utils"
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

func TestBlackjack_double_0(t *testing.T) {
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

func TestBlackjack_double_1(t *testing.T) {
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

	err := testBlackjack.double(&testPlayer, 1)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, uint64(10), testPlayer.Hands[0].BetAmount, "bet amount unchanged")
	require.False(t, testPlayer.Hands[0].isDoubled, "is not doubled")

	require.Equal(t, uint64(100), testPlayer.Hands[1].BetAmount, "bet amount doubled")
	require.True(t, testPlayer.Hands[1].isDoubled, "is doubled")
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
	testBlackjack.deck.BurntCards = []deck.Card{}

	err := testBlackjack.double(&testPlayer, 0)

	require.EqualError(t, err, "Failed to dealPlayerACard: Failed to DealCard: Failed to TakeCardByIndex: Index is out of range")
}

func Test_autoDrawCards(t *testing.T) {
	type testCase struct {
		name                string
		inputHand           Hand
		inputDeck           deck.Deck
		shouldHaveBlackjack bool
		shouldDraw          bool
		expectedError       error
	}

	testCases := []testCase{
		{
			name: "HasBlackjack",
			inputHand: Hand{
				cards: []deck.Card{
					deck.AceOfClubs,
				},
			},
			inputDeck: deck.Deck{
				ActiveCards: []deck.Card{
					deck.TenOfClubs,
				},
			},
			shouldHaveBlackjack: true,
			shouldDraw:          true,
		},
		{
			name: "DoesNotHaveBlackjack",
			inputHand: Hand{
				cards: []deck.Card{
					deck.AceOfClubs,
				},
			},
			inputDeck: deck.Deck{
				ActiveCards: []deck.Card{
					deck.NineOfClubs,
				},
			},
			shouldHaveBlackjack: false,
			shouldDraw:          true,
		},
		{
			name: "HasNoChanceOfBlackjack",
			inputHand: Hand{
				cards: []deck.Card{
					deck.NineOfHearts,
				},
			},
			inputDeck: deck.Deck{
				ActiveCards: []deck.Card{
					deck.TenOfClubs,
				},
			},
			shouldHaveBlackjack: false,
			shouldDraw:          false,
		},
	}

	testPlayer := player.Player{}
	testBJPlayer := NewBlackjackPlayer(testPlayer, nil)
	testBJPlayers := []BlackjackPlayer{testBJPlayer}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			testBlackjack.deck = &tc.inputDeck
			testBJPlayers[0].Hands = []Hand{tc.inputHand}

			actualErr := testBlackjack.autoDrawCards(testBJPlayers)

			if tc.expectedError != nil {
				require.EqualError(t, actualErr, tc.expectedError.Error())
				return
			}

			actualHand := testBJPlayers[0].Hands[0]

			require.NoError(t, actualErr)
			require.Equal(t, tc.shouldHaveBlackjack, actualHand.Blackjack())

			if tc.shouldDraw {
				require.Len(t, actualHand.cards, 2)
			} else {
				require.Len(t, actualHand.cards, 1)
			}
		})
	}
}

func TestBlackjack_handlePotentialSplits(t *testing.T) {
	type testCase struct {
		name            string
		inputPlayer     *BlackjackPlayer
		inputDealerHand DealerHand
		expectedErr     error
	}

	testCases := []testCase{
		{
			name: "Error InvalidBetAmount",
			inputPlayer: &BlackjackPlayer{
				Hands: []Hand{
					{
						BetAmount: 100,
					},
				},
				Player: player.Player{
					Credits: 99,
				},
			},
			inputDealerHand: DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.AceOfClubs,
					},
				},
			},
			expectedErr: errors.New("Insufficient credits"),
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			err := testBlackjack.handlePotentialSplits(tc.inputPlayer, tc.inputDealerHand)

			if tc.expectedErr != nil {
				require.EqualError(t, tc.expectedErr, err.Error(), "error value")
				return
			}

			require.NoError(t, err, "unexpected error")
		})
	}
}

func Test_ShouldContinuePlaying(t *testing.T) {
	type testCase struct {
		name        string
		inputAction string
		inputHand   Hand
		expectedRes bool
	}

	testCases := []testCase{
		{
			name:        "Stand",
			inputAction: bjutils.Stand,
			inputHand: Hand{
				cards: []deck.Card{
					deck.ThreeOfClubs,
				},
			},
			expectedRes: false,
		},
		{
			name:        "Bust",
			inputAction: bjutils.Hit,
			inputHand: Hand{
				cards: []deck.Card{
					deck.TenOfDiamonds,
					deck.KingOfSpades,
					deck.EightOfClubs,
				},
			},
			expectedRes: false,
		},
		{
			name:        "21",
			inputAction: bjutils.Hit,
			inputHand: Hand{
				cards: []deck.Card{
					deck.EightOfHearts,
					deck.TwoOfSpades,
					deck.ThreeOfSpades,
					deck.FourOfDiamonds,
					deck.FourOfDiamonds,
				},
			},
			expectedRes: false,
		},
		{
			name:        "ShouldContinue",
			inputAction: bjutils.Hit,
			inputHand: Hand{
				cards: []deck.Card{
					deck.ThreeOfHearts,
				},
			},
			expectedRes: true,
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			actualRes := ShouldContinuePlaying(tc.inputAction, tc.inputHand)

			require.Equal(t, tc.expectedRes, actualRes)
		})
	}
}
