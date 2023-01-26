package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) Play(logger *zap.Logger, players []BlackJackPlayer, dealerHand DealerHand) error {
	logger.Info("playing round")
	for i, p := range players {
		for j, h := range p.Hands {
			logger.Info("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

			input := "y"

			dealerHand.DealerLog(logger)
			for input == "y" && !h.Bust() {
				h.Log(logger)
				fmt.Println("Take card? (y/N)")
				_, err := fmt.Scanln(&input)
				if err != nil {
					return errors.ErrFailedSubMethod("fmt.Scanln", err)
				}

				if input == "y" {
					var c *deck.Card
					c, err = bj.DealCard()
					if err != nil {
						return errors.ErrFailedSubMethod("DealCard", err)
					}

					h.AddCard(*c)
					c.Log(logger)
				}
			}

			if h.Bust() {
				logger.Info("player bust")
			}
		}
	}

	if err := bj.Results(logger, players, dealerHand); err != nil {
		return errors.ErrFailedSubMethod("Results", err)
	}

	for _, p := range players {
		if err := p.PrintResult(logger); err != nil {
			return errors.ErrFailedSubMethod("PrintResult", err)
		}
	}

	return nil
}
