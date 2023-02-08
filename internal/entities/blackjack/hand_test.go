package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	internalTesting "scrub/internal/testing"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHand_CanSplit(t *testing.T) {
	type test struct {
		name     string
		input    Hand
		expected bool
	}

	testCases := []test{
		{
			name: "Can",
			input: Hand{
				cards: []deck.Card{
					{
						Symbol: "10",
						Suit:   "Hearts",
						Value:  10,
					},
					{
						Symbol: "10",
						Suit:   "Diamonds",
						Value:  10,
					},
				},
			},
			expected: true,
		},
		{
			name: "Cannot",
			input: Hand{
				cards: []deck.Card{
					{
						Symbol: "10",
						Suit:   "Hearts",
						Value:  10,
					},
					{
						Symbol: "9",
						Suit:   "Diamonds",
						Value:  9,
					},
				},
			},
			expected: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			actual := tc.input.CanSplit()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestHand_IsSoft(t *testing.T) {
	type test struct {
		name     string
		input    Hand
		expected bool
	}

	testCases := []test{
		{
			name: "Soft",
			input: Hand{
				cards: []deck.Card{
					{
						Symbol: "Ace",
						Suit:   "Hearts",
						Value:  1,
					},
					{
						Symbol: "10",
						Suit:   "Diamonds",
						Value:  10,
					},
				},
			},
			expected: true,
		},
		{
			name: "Soft",
			input: Hand{
				cards: []deck.Card{
					{
						Symbol: "Ace",
						Suit:   "Hearts",
						Value:  1,
					},
					{
						Symbol: "6",
						Suit:   "Diamonds",
						Value:  6,
					},
				},
			},
			expected: true,
		},
		{
			name: "Hard",
			input: Hand{
				cards: []deck.Card{
					{
						Symbol: "9",
						Suit:   "Hearts",
						Value:  9,
					},
					{
						Symbol: "10",
						Suit:   "Diamonds",
						Value:  10,
					},
				},
			},
			expected: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			actual := tc.input.IsSoft()
			require.Equal(t, tc.expected, actual)
		})
	}
}
