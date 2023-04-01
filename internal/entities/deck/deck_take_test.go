package deck

import (
	"scrub/internal/errorutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeck_TakeCardByIndex(t *testing.T) {
	type testCase struct {
		name          string
		indexes       []int
		expectedCards []*Card
		expectedError error
	}

	testCases := []testCase{
		{
			name:    "OK First Card",
			indexes: []int{0},
			expectedCards: []*Card{
				&AceOfClubs,
			},
			expectedError: nil,
		},
		{
			name:    "OK Last Card",
			indexes: []int{51},
			expectedCards: []*Card{
				&KingOfSpades,
			},
			expectedError: nil,
		},
		{
			name:          "Err Index Out Of Range",
			indexes:       []int{52},
			expectedError: errorutils.ErrIndexOutOfRange,
		},
		{
			name:    "OK Multiple Cards",
			indexes: []int{0, 0},
			expectedCards: []*Card{
				&AceOfClubs,
				&TwoOfClubs,
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDeck := NewDeck()

			var err error
			cards := make([]*Card, len(tc.indexes))
			for i, index := range tc.indexes {
				cards[i], err = testDeck.TakeCardByIndex(index)
				if err != nil {
					break
				}
			}

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedError, err, "error message")
				return
			}

			require.NoError(t, err, "must not return error")
			assert.EqualValues(t, cards, tc.expectedCards, "selected cards")
			assert.Len(t, testDeck.ActiveCards, 52-len(tc.indexes), "%d cards should remain", 52-len(tc.indexes))
			assert.Len(t, testDeck.BurntCards, len(tc.indexes), "%d card should be burnt", len(tc.indexes))
		})
	}
}
