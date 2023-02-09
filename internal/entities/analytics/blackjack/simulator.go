package blackjackanalytics

import (
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/player"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func Simulate(rounds, decks uint) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	players := []blackjack.BlackJackPlayer{
		{
			Player: player.Player{
				Name:    "Test Player",
				Credits: 1000000,
			},
			Hands: []blackjack.Hand{
				{
					BetAmount: 10,
				},
			},
		},
	}

	bj := blackjack.NewBlackjack(decks)

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

	logger.Info("simulation complete", zap.Uint("rounds", i), zap.Uint64("credits", players[0].Credits))
	return nil
}
