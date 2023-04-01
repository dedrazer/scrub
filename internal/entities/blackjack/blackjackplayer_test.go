package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/testutils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlackJackPlayer_ResetHands(t *testing.T) {
	type testCase struct {
		name           string
		input          *BlackjackPlayer
		expectedOutput *BlackjackPlayer
		expectedError  error
	}

	testCases := []testCase{
		{
			name: "Doubled",
			input: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.AceOfClubs,
							deck.TenOfDiamonds,
						},
						BetAmount: 100,
						isDoubled: true,
					},
				},
			},
			expectedOutput: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.AceOfClubs,
							deck.TenOfDiamonds,
						},
						BetAmount: 50,
						isDoubled: false,
					},
				},
			},
		},
		{
			name: "Split",
			input: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.FourOfHearts,
							deck.FiveOfDiamonds,
						},
						BetAmount: 50,
						isSplit:   true,
					},
					{
						cards: []deck.Card{
							deck.FourOfClubs,
							deck.SixOfDiamonds,
						},
						BetAmount: 50,
						isSplit:   true,
					},
				},
			},
			expectedOutput: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.FourOfClubs,
							deck.SixOfDiamonds,
						},
						BetAmount: 50,
						isSplit:   false,
					},
				},
			},
		},
		{
			name: "Split & Doubled",
			input: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.FourOfHearts,
							deck.FiveOfDiamonds,
						},
						BetAmount: 50,
						isSplit:   true,
					},
					{
						cards: []deck.Card{
							deck.FourOfClubs,
							deck.SixOfDiamonds,
						},
						BetAmount: 100,
						isSplit:   true,
						isDoubled: true,
					},
				},
			},
			expectedOutput: &BlackjackPlayer{
				Hands: []Hand{
					{
						cards: []deck.Card{
							deck.FourOfClubs,
							deck.SixOfDiamonds,
						},
						BetAmount: 50,
						isSplit:   false,
					},
				},
			},
		},
	}

	for testNumber, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, testNumber, tc.name), func(t *testing.T) {
			err := tc.input.ResetHands()
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedOutput, tc.input)
		})
	}
}
