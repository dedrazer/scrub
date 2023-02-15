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
					deck.TenOfHearts,
					deck.TenOfDiamonds,
				},
			},
			expected: true,
		},
		{
			name: "Cannot",
			input: Hand{
				cards: []deck.Card{
					deck.TenOfHearts,
					deck.NineOfDiamonds,
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
					deck.AceOfHearts,
					deck.TenOfDiamonds,
				},
			},
			expected: true,
		},
		{
			name: "Soft",
			input: Hand{
				cards: []deck.Card{
					deck.AceOfHearts,
					deck.SixOfDiamonds,
				},
			},
			expected: true,
		},
		{
			name: "Hard",
			input: Hand{
				cards: []deck.Card{
					deck.NineOfHearts,
					deck.TenOfDiamonds,
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
