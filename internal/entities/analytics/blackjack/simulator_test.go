package blackjackanalytics

import (
	"os"
	"scrub/internal/entities/blackjack/bettingstrategy"
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

func TestSimulator_getOneCreditPercentageOfTotal(t *testing.T) {
	testSimulator.OneCreditAmount = 9
	testSimulator.StartingCredits = 100
	actual := testSimulator.getOneCreditPercentageOfTotal()

	require.Equal(t, float64(0.09), actual)
}
