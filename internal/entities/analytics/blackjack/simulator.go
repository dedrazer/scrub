package blackjackanalytics

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/player"
	"scrub/internal/errors"
	"time"

	"go.uber.org/zap"
)

type SimulationConfig struct {
	Rounds          uint
	Decks           uint
	StartingCredits uint64
	OneCreditAmount uint64
	RebuyCount      int
}

func Simulate(logger *zap.Logger, simulationConfig SimulationConfig, bettingStrategy func(*zap.Logger, []blackjack.BlackjackPlayer, uint64) error) error {
	logger.Info("starting simulation", zap.Uint("rounds", simulationConfig.Rounds), zap.Uint("decks", simulationConfig.Decks))

	players := []blackjack.BlackjackPlayer{
		{
			Player: player.Player{
				Name:            "Test Player",
				StartingCredits: simulationConfig.StartingCredits,
				Credits:         simulationConfig.StartingCredits,
			},
			Hands: []blackjack.Hand{
				{
					BetAmount: simulationConfig.OneCreditAmount,
				},
			},
		},
	}

	bj := blackjack.NewBlackjack(simulationConfig.Decks)

	startTime := time.Now().UTC()
	startingRebuyCount := simulationConfig.RebuyCount

	var i uint = 0
	for i < simulationConfig.Rounds && (players[0].Credits > 0 || simulationConfig.RebuyCount > 0) {
		i++

		if players[0].Credits == 0 {
			logger.Info("player rebuy", zap.Int("rebuys remaining", simulationConfig.RebuyCount))
			players[0].Credits = simulationConfig.StartingCredits
			simulationConfig.RebuyCount--
		}

		err := bettingStrategy(logger, players, simulationConfig.OneCreditAmount)
		if err != nil {
			return errors.ErrFailedSubMethod("bettingStrategy", err)
		}

		dealerHand, err := bj.DealRound(players)
		if err != nil {
			return errors.ErrFailedSubMethod("DealRound", err)
		}

		err = bj.Play(logger, players, dealerHand, blackjack.Strategy)
		if err != nil {
			return errors.ErrFailedSubMethod("Play", err)
		}
	}

	logger.Info("simulation complete", zap.Uint("rounds", i), zap.Int("rebuys", startingRebuyCount-simulationConfig.RebuyCount))

	for j := range players {
		players[j].LogStatistics(logger)
	}

	bj.LogStatistics(logger)

	totalDurationMs := time.Since(startTime).Milliseconds()
	totalDuration := fmt.Sprintf("%dms", totalDurationMs)
	averageRoundDuration := fmt.Sprintf("%.2fμs", (float64(totalDurationMs)/float64(simulationConfig.Rounds))*1000)
	roundsPerSecond := int64(float64(int64(simulationConfig.Rounds)*1000) / float64(totalDurationMs))
	logger.Info("runtime statistics",
		zap.String("duration", totalDuration),
		zap.Int64("rounds per second", roundsPerSecond),
		zap.String("average round duration", averageRoundDuration))
	return nil
}
