package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) Play(logger *zap.Logger, playerHands []Hand, dealerHand DealerHand) error {
	logger.Info("playing round")
	for i, h := range playerHands {
		logger.Info("turn", zap.Int("player", i+1))

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

	results, err := bj.Results(logger, playerHands, dealerHand)
	if err != nil {
		return errors.ErrFailedSubMethod("Results", err)
	}

	for i, result := range results {
		player := "player " + fmt.Sprint(i+1)
		status := ""

		if result == nil {
			status = "push"
		} else if *result == false {
			status = "lose"
		} else if *result == true {
			status = "win"
		}
		logger.Info(player, zap.String("status", status))
	}

	return nil
}
