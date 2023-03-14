package blackjackanalytics

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/bettingstrategy"
	"scrub/internal/entities/blackjack/models"
	"scrub/internal/entities/player"
	"scrub/internal/errors"
	"time"

	"go.uber.org/zap"
)

type SimulationConfig struct {
	Rounds          uint
	Decks           uint
	StartingCredits uint64
	BankCredits     uint64
	BankAtCredits   uint64
	OneCreditAmount uint64
	RebuyCount      int
}

func Simulate(logger *zap.Logger, simulationConfig SimulationConfig, bettingStrategy bettingstrategy.Strategy) error {
	logger.Info("starting simulation", zap.Any("config", simulationConfig))

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
	startingBankCredits := simulationConfig.BankCredits

	lastCreditRound := uint(0)
	creditAtRound := make([]uint, 1)
	highestProfitPercentage := float64(1)

	var numberOfDeposits, numberOfWithdrawals uint

	var i uint = 0
	for i < simulationConfig.Rounds && (players[0].Credits > 0 || (simulationConfig.RebuyCount > 0 && simulationConfig.BankCredits > 0)) {
		logger.Debug("starting round", zap.Uint("round", i))

		if players[0].Credits == 0 {
			roundsSinceLastCredit := i - lastCreditRound

			if simulationConfig.BankCredits >= simulationConfig.StartingCredits {
				players[0].Credits = simulationConfig.StartingCredits
				simulationConfig.BankCredits -= simulationConfig.StartingCredits
				numberOfWithdrawals++
				logger.Debug("withdrew", zap.Uint64("credits", simulationConfig.StartingCredits), zap.Uint("round", i))
			} else {
				players[0].Credits = simulationConfig.BankCredits
				simulationConfig.BankCredits = 0
				logger.Info("bank is out of credits", zap.Uint("round", i))
				numberOfWithdrawals++
				logger.Debug("withdrew", zap.Uint64("credits", players[0].Credits), zap.Uint("round", i))
			}

			players[0].Hands[0].BetAmount = simulationConfig.OneCreditAmount
			simulationConfig.RebuyCount--
			logger.Debug("player rebuy", zap.Int("rebuys remaining", simulationConfig.RebuyCount))

			lastCreditRound = i
			creditAtRound = append(creditAtRound, roundsSinceLastCredit)
		}

		if players[0].Credits >= simulationConfig.BankAtCredits {
			simulationConfig.BankCredits += players[0].Credits - simulationConfig.StartingCredits
			numberOfDeposits++
			logger.Debug("deposited", zap.Uint64("credits", players[0].Credits-simulationConfig.StartingCredits), zap.Uint("round", i))
			players[0].Credits = simulationConfig.StartingCredits
		}

		profitPercentage := float64(players[0].Credits) / float64(simulationConfig.StartingCredits)
		if profitPercentage > highestProfitPercentage {
			highestProfitPercentage = profitPercentage
		}

		err := bettingStrategy.Strategy(players)
		if err != nil {
			return errors.ErrFailedSubMethod("bettingStrategy", err)
		}

		if players[0].Hands[0].BetAmount > players[0].Credits {
			players[0].Hands[0].BetAmount = players[0].Credits
		}

		dealerHand, err := bj.DealRound(players)
		if err != nil {
			return errors.ErrFailedSubMethod("DealRound", err)
		}

		err = bj.Play(logger, players, dealerHand, blackjack.Strategy)
		if err != nil {
			return errors.ErrFailedSubMethod("Play", err)
		}

		i++
	}

	depositPercentage := float64(numberOfDeposits) / float64(numberOfWithdrawals)

	logger.Info("simulation complete",
		zap.Uint("rounds", i),
		zap.Int("remaining rebuys", startingRebuyCount-simulationConfig.RebuyCount),
		zap.Uint("deposits", numberOfDeposits),
		zap.Uint("withdrawals", numberOfWithdrawals),
		zap.String("deposit percentage", fmt.Sprintf("%.2f%%", depositPercentage*100)))

	for j := range players {
		players[j].LogStatistics(logger)
	}

	bj.LogStatistics(logger)

	totalDurationMs := time.Since(startTime).Milliseconds()
	var totalDurationTextual string

	if totalDurationMs > 10000 {
		totalDurationTextual = fmt.Sprintf("%ds", totalDurationMs/1000)
	} else {
		totalDurationTextual = fmt.Sprintf("%dms", totalDurationMs)
	}

	averageRoundDuration := fmt.Sprintf("%.2fμs", (float64(totalDurationMs)/float64(simulationConfig.Rounds))*1000)
	roundsPerSecond := int64(float64(int64(simulationConfig.Rounds)*1000) / float64(totalDurationMs))
	logger.Info("runtime statistics",
		zap.String("duration", totalDurationTextual),
		zap.Int64("rounds per second", roundsPerSecond),
		zap.String("average round duration", averageRoundDuration))

	totalRounds := uint64(0)
	earliestBankruptcyRound := i
	for j := range creditAtRound {
		if creditAtRound[j] != 0 && creditAtRound[j] < earliestBankruptcyRound {
			earliestBankruptcyRound = creditAtRound[j]
		}

		totalRounds += uint64(creditAtRound[j])
	}
	averageRoundsSurvived := float64(totalRounds) / float64(len(creditAtRound))
	oneCreditPercentageOfTotal := float64(simulationConfig.OneCreditAmount) / float64(simulationConfig.StartingCredits)

	res := models.SimulationResults{
		AverageRoundsSurvived:      uint(averageRoundsSurvived),
		EarliestBankruptcyRound:    earliestBankruptcyRound,
		HighestProfitPercentage:    highestProfitPercentage,
		OneCreditPercentageOfTotal: oneCreditPercentageOfTotal,
		StartingCredits:            startingBankCredits,
		EndingCredits:              simulationConfig.BankCredits,
		RebuyCredits:               simulationConfig.StartingCredits,
		BankAtCredits:              simulationConfig.BankAtCredits,
		Score:                      float64(highestProfitPercentage) * depositPercentage * float64(averageRoundsSurvived) * float64(oneCreditPercentageOfTotal),
	}

	logger.Info("strategy results", zap.Any("results", res))

	return nil
}
