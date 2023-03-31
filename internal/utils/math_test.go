package utils

import (
	"fmt"
	"testing"

	"github.com/fasttrack-solutions/altenar-transformer/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestRound(t *testing.T) {
	type testCase struct {
		name               string
		inputValue         float64
		inputDecimalPlaces int
		expectedResult     float64
	}

	testCases := []testCase{
		{
			name:               "Round down 2dp",
			inputValue:         1.234567,
			inputDecimalPlaces: 2,
			expectedResult:     1.23,
		},
		{
			name:               "Round up 2dp",
			inputValue:         1.235567,
			inputDecimalPlaces: 2,
			expectedResult:     1.24,
		},
		{
			name:               "Round down 0dp",
			inputValue:         1.234567,
			inputDecimalPlaces: 0,
			expectedResult:     1,
		},
		{
			name:               "Round down 6dp",
			inputValue:         1.2345671,
			inputDecimalPlaces: 6,
			expectedResult:     1.234567,
		},
		{
			name:               "Round up 6dp",
			inputValue:         1.2345675,
			inputDecimalPlaces: 6,
			expectedResult:     1.234568,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf(testutils.TestNameTemplate, i, tc.name), func(t *testing.T) {
			result := Round(tc.inputValue, tc.inputDecimalPlaces)

			require.Equal(t, tc.expectedResult, result)
		})
	}
}
