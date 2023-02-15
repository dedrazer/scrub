package blackjackanalytics

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/player"
	"scrub/internal/errors"
	"time"

	"go.uber.org/zap"
)

func Simulate(rounds, decks uint) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	startingCredits := uint64(1000000)

	players := []blackjack.BlackJackPlayer{
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

		err = bj.Play(logger, players, dealerHand, blackjack.Strategy1)
		if err != nil {
			return errors.ErrFailedSubMethod("Play", err)
		}
	}

	logger.Info("simulation complete", zap.Uint("rounds", i))

	for j := range players {
		players[j].LogStatistics(logger)
	}

	totalDurationMs := time.Since(startTime).Milliseconds()
	totalDuration := fmt.Sprintf("%dms", totalDurationMs)
	averageRoundDuration := fmt.Sprintf("%.2fμs", (float64(totalDurationMs)/float64(rounds))*1000)
	logger.Info("runtime statistics",
		zap.String("duration", totalDuration),
		zap.String("average round duration", averageRoundDuration))
	return nil
}
