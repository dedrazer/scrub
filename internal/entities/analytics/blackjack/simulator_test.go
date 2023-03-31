package blackjackanalytics

import (
	"fmt"
	"os"
	"scrub/internal/entities/analytics/blackjack/models"
	"scrub/internal/entities/blackjack/bettingstrategy"
	internalTesting "scrub/internal/testing"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testSimulator *Simulator
)

func TestMain(m *testing.M) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	testSimulationConfig := SimulationConfig{
		MaxRounds:       100,
		Decks:           20,
		BankCredits:     1000,
		BankAtCredits:   100,
		StartingCredits: 50,
		OneCreditAmount: 10,
		RebuyCount:      100,
	}

	testStrategy := &bettingstrategy.Martingale{
		CommonStrategyVariables: bettingstrategy.CommonStrategyVariables{
			OneCreditValue: testSimulationConfig.OneCreditAmount,
			Logger:         logger,
		},
	}

	testSimulator = NewSimulator(logger, testStrategy, testSimulationConfig)

	runCode := m.Run()

	os.Exit(runCode)
}

func TestSimulator_getSimulationResults(t *testing.T) {
	testSimulator.OneCreditAmount = 10
	testSimulator.startingBankCredits = 1000
	testSimulator.StartingCredits = 100
	testSimulator.averageRoundsSurvived = 31
	testSimulator.earliestBankruptcyRound = 5
	testSimulator.highestProfitPercentage = 3.14
	testSimulator.BankCredits = 800
	testSimulator.BankAtCredits = 100

	expected := models.SimulationResults{
		AverageRoundsSurvived:      31,
		EarliestBankruptcyRound:    5,
		HighestProfitPercentage:    3.14,
		OneCreditPercentageOfTotal: 0.1,
		StartingCredits:            1000,
		EndingCredits:              800,
		RebuyCredits:               100,
		BankAtCredits:              100,
		Score:                      9.73,
	}

	actual := testSimulator.getSimulationResults()

	require.Equal(t, expected, actual)
}

func TestSimulator_getTextualDuration(t *testing.T) {
	type testCase struct {
		name            string
		inputDurationMs int64
		expected        string
	}

	testCases := []testCase{
		{
			name:            "Zero",
			inputDurationMs: 0,
			expected:        "0ms",
		},
		{
			name:            "Seconds",
			inputDurationMs: 1000,
			expected:        "1.00sec",
		},
		{
			name:            "Minutes",
			inputDurationMs: 60000,
			expected:        "1.00min",
		},
		{
			name:            "Hours",
			inputDurationMs: 3600000,
			expected:        "1.00hrs",
		},
		{
			name:            "FractionalSeconds",
			inputDurationMs: 1234,
			expected:        "1.23sec",
		},
		{
			name:            "FractionalMinutes",
			inputDurationMs: 123456,
			expected:        "2.06min",
		},
		{
			name:            "FractionalHours",
			inputDurationMs: 12345678,
			expected:        "3.43hrs",
		},
		{
			name:            "Negative",
			inputDurationMs: -1000,
			expected:        "-1000ms",
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			actual := testSimulator.getTextualDuration(tc.inputDurationMs)

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSimulator_getScore(t *testing.T) {
	testSimulator.highestProfitPercentage = 3.14
	testSimulator.numberOfDeposits = 56
	testSimulator.numberOfWithdrawals = 13
	testSimulator.averageRoundsSurvived = 21
	testSimulator.OneCreditAmount = 50
	testSimulator.StartingCredits = 100

	actual := testSimulator.getScore()
	require.Equal(t, 142.02, actual)
}

func TestSimulator_getOneCreditPercentageOfStartingCredits(t *testing.T) {
	testSimulator.OneCreditAmount = 9
	testSimulator.StartingCredits = 100
	actual := testSimulator.getOneCreditPercentageOfStartingCredits()

	require.Equal(t, float64(0.09), actual)
}

func TestSimulator_getDepositPercentage(t *testing.T) {
	type testCase struct {
		name                     string
		inputNumberOfDeposits    uint
		inputNumberOfWithdrawals uint
		expected                 float64
	}

	testCases := []testCase{
		{
			name:                     "ZeroWithdrawals DivByZeroCheck",
			inputNumberOfDeposits:    2,
			inputNumberOfWithdrawals: 0,
			expected:                 1,
		},
		{
			name:                     "ZeroDeposits",
			inputNumberOfDeposits:    0,
			inputNumberOfWithdrawals: 1,
			expected:                 0,
		},
		{
			name:                     "Normal Loss",
			inputNumberOfDeposits:    2,
			inputNumberOfWithdrawals: 6,
			expected:                 float64(2) / 6,
		},
		{
			name:                     "Normal Profit",
			inputNumberOfDeposits:    5,
			inputNumberOfWithdrawals: 3,
			expected:                 float64(5) / 3,
		},
	}

	for tn, tc := range testCases {
		t.Run(fmt.Sprintf(internalTesting.TestNameTemplate, tn, tc.name), func(t *testing.T) {
			testSimulator.numberOfDeposits = tc.inputNumberOfDeposits
			testSimulator.numberOfWithdrawals = tc.inputNumberOfWithdrawals

			actual := testSimulator.getDepositPercentage()

			require.Equal(t, tc.expected, actual)
		})
	}
}
