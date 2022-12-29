package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) Play(logger *zap.Logger, playerHands []Hand, dealerHand Hand) error {
	logger.Info("playing round")
	for i, h := range playerHands {
		logger.Info("turn", zap.Int("player", i+1))

		input := "y"

		for input == "y" && !h.Bust() {
			dealerHand.DealerLog(logger)
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
	}

	results, err := bj.Results(logger, playerHands, dealerHand)
	if err != nil {
		return errors.ErrFailedSubMethod("Results", err)
	}

	for i, result := range results {
		logger.Info("result", zap.Int("player", i+1), zap.Any("result", result))
	}

	return nil
}
