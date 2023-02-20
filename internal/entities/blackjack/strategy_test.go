package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	internalTesting "scrub/internal/testing"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrategy1(t *testing.T) {
	type test struct {
		name          string
		playerSymbols []string
		dealerSymbol  string
		expected      string
	}

	testCases := []test{
		{
			name:          "player should hit on soft 17",
			playerSymbols: []string{"Ace", "6"},
			dealerSymbol:  "7",
			expected:      hit,
		},
		{
			name:          "player should stand on hard 17",
			playerSymbols: []string{"7", "10"},
			dealerSymbol:  "6",
			expected:      stand,
		},
		{
			name:          "player should hit on soft 18",
			playerSymbols: []string{"Ace", "7"},
			dealerSymbol:  "9",
			expected:      hit,
		},
		{
			name:          "player should stand on hard 18",
			playerSymbols: []string{"8", "10"},
			dealerSymbol:  "Q",
			expected:      stand,
		},
		{
			name:          "player should stand on 19",
			playerSymbols: []string{"10", "9"},
			dealerSymbol:  "5",
			expected:      stand,
		},
		{
			name:          "player should hit on 12",
			playerSymbols: []string{"4", "8"},
			dealerSymbol:  "Q",
			expected:      hit,
		},
		{
			name:          "player should stand on 20",
			playerSymbols: []string{"10", "10"},
			dealerSymbol:  "2",
			expected:      stand,
		},
		{
			name:          "player double down on hard 11",
			playerSymbols: []string{"6", "5"},
			dealerSymbol:  "J",
			expected:      double,
		},
		{
			name:          "player should stand on soft 21",
			playerSymbols: []string{"10", "Ace"},
			dealerSymbol:  "J",
			expected:      stand,
		},
		{
			name:          "player should double down on 10",
			playerSymbols: []string{"5", "5"},
			dealerSymbol:  "4",
			expected:      double,
		},
		{
			name:          "player should double down on 9",
			playerSymbols: []string{"4", "5"},
			dealerSymbol:  "3",
			expected:      double,
		},
	}

	for testNumber, testCase := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, testNumber, testCase.name), func(t *testing.T) {
			dealerHand := DealerHand{
				Hand: Hand{
					cards: []deck.Card{
						deck.TenOfHearts,
						{Symbol: testCase.dealerSymbol, Value: deck.CardValues[testCase.dealerSymbol]},
					},
				},
			}
			playerHand := Hand{cards: []deck.Card{}}
			for _, symbol := range testCase.playerSymbols {
				playerHand.cards = append(playerHand.cards, deck.Card{Symbol: symbol, Value: deck.CardValues[symbol]})
			}

			actual := Strategy(playerHand, dealerHand)
			require.Equal(t, testCase.expected, actual)
		})
	}
}
