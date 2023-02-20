package blackjackanalytics

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/player"
	"scrub/internal/errors"
	"time"

	"go.uber.org/zap"
)

func Simulate(logger *zap.Logger, rounds, decks uint) error {
	logger.Info("starting simulation", zap.Uint("rounds", rounds), zap.Uint("decks", decks))

	startingCredits := uint64(1000000)

	players := []blackjack.BlackjackPlayer{
		{
			Player: player.Player{
				Name:            "Test Player",
				StartingCredits: startingCredits,
				Credits:         startingCredits,
			},
			Hands: []blackjack.Hand{
				{
					BetAmount: 10,
				},
			},
		},
	}

	bj := blackjack.NewBlackjack(decks)

	startTime := time.Now().UTC()

	var i uint = 0
	for i < rounds && players[0].Credits > 0 {
		i++

		dealerHand, err := bj.DealRound(players)
		if err != nil {
			return errors.ErrFailedSubMethod("DealRound", err)
		}

		err = bj.Play(logger, players, dealerHand, blackjack.Strategy)
		if err != nil {
			return errors.ErrFailedSubMethod("Play", err)
		}
	}

	logger.Info("simulation complete", zap.Uint("rounds", i))

	for j := range players {
		players[j].LogStatistics(logger)
	}

	bj.LogStatistics(logger)

	totalDurationMs := time.Since(startTime).Milliseconds()
	totalDuration := fmt.Sprintf("%dms", totalDurationMs)
	averageRoundDuration := fmt.Sprintf("%.2fμs", (float64(totalDurationMs)/float64(rounds))*1000)
	roundsPerSecond := int64(float64(int64(rounds)*1000) / float64(totalDurationMs))
	logger.Info("runtime statistics",
		zap.String("duration", totalDuration),
		zap.Int64("rounds per second", roundsPerSecond),
		zap.String("average round duration", averageRoundDuration))
	return nil
}
