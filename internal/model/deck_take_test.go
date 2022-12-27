package model

import (
	"scrub/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck_TakeCardByIndex(t *testing.T) {
	type testCase struct {
		name          string
		index         int
		expectedCard  string
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "OK First Card",
			index:         0,
			expectedCard:  "Ace of Clubs",
			expectedError: nil,
		},
		{
			name:          "OK Last Card",
			index:         51,
			expectedCard:  "King of Spades",
			expectedError: nil,
		},
		{
			name:          "Err Index Out Of Range",
			index:         52,
			expectedError: errors.ErrIndexOutOfRange,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDeck := NewDeck()

			card, err := testDeck.TakeCardByIndex(tc.index)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedError, err, "error message")
				return
			}

			require.NoError(t, err, "must not return error")
			assert.Equal(t, card.Print(), tc.expectedCard, "selected card")
			assert.Len(t, testDeck.ActiveCards, 51, "51 cards should remain")
			assert.Len(t, testDeck.BurntCards, 1, "1 card should be burnt")
		})
	}
}
