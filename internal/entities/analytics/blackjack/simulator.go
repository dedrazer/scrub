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

type Simulator struct {
	logger                  *zap.Logger
	strategy                bettingstrategy.Strategy
	startingTime            time.Time
	startingRebuyCount      int
	startingBankCredits     uint64
	lastCreditRound         uint
	creditAtRound           []uint
	highestProfitPercentage float64
	numberOfDeposits        uint
	numberOfWithdrawals     uint
	currentRound            uint
	SimulationConfig
}

func NewSimulator(logger *zap.Logger, strategy bettingstrategy.Strategy, config SimulationConfig) *Simulator {
	return &Simulator{
		logger:                  logger,
		strategy:                strategy,
		startingTime:            time.Now().UTC(),
		startingRebuyCount:      config.RebuyCount,
		startingBankCredits:     config.BankCredits,
		lastCreditRound:         uint(0),
		creditAtRound:           make([]uint, 1),
		highestProfitPercentage: float64(1),
	}
}

func (s *Simulator) Simulate(logger *zap.Logger, simulationConfig SimulationConfig, bettingStrategy bettingstrategy.Strategy) error {
	logger.Info("starting simulation", zap.Any("config", simulationConfig))

	players := initTestPlayers(simulationConfig)

	bj := blackjack.NewBlackjack(simulationConfig.Decks)

	s.currentRound = 0
	for s.currentRound < simulationConfig.Rounds && (players[0].Credits > 0 || (simulationConfig.RebuyCount > 0 && simulationConfig.BankCredits > 0)) {
		err := s.simulationRound(logger, simulationConfig, bettingStrategy, players, bj, s.currentRound)
		if err != nil {
			return errors.ErrFailedSubMethod("simulationRound", err)
		}

		s.currentRound++
	}

	s.logSimulationCompletion()

	for j := range players {
		players[j].LogStatistics(logger)
	}

	bj.LogStatistics(logger)

	totalDurationMs := time.Since(s.startingTime).Milliseconds()
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
	earliestBankruptcyRound := s.currentRound
	for j := range s.creditAtRound {
		if s.creditAtRound[j] != 0 && s.creditAtRound[j] < earliestBankruptcyRound {
			earliestBankruptcyRound = s.creditAtRound[j]
		}

		totalRounds += uint64(s.creditAtRound[j])
	}
	averageRoundsSurvived := float64(totalRounds) / float64(len(s.creditAtRound))
	oneCreditPercentageOfTotal := float64(simulationConfig.OneCreditAmount) / float64(simulationConfig.StartingCredits)

	res := models.SimulationResults{
		AverageRoundsSurvived:      uint(averageRoundsSurvived),
		EarliestBankruptcyRound:    earliestBankruptcyRound,
		HighestProfitPercentage:    s.highestProfitPercentage,
		OneCreditPercentageOfTotal: oneCreditPercentageOfTotal,
		StartingCredits:            s.startingBankCredits,
		EndingCredits:              simulationConfig.BankCredits,
		RebuyCredits:               simulationConfig.StartingCredits,
		BankAtCredits:              simulationConfig.BankAtCredits,
		Score:                      float64(s.highestProfitPercentage) * s.getDepositPercentage() * float64(averageRoundsSurvived) * float64(oneCreditPercentageOfTotal),
	}

	logger.Info("strategy results", zap.Any("results", res))

	return nil
}

func (s *Simulator) simulationRound(logger *zap.Logger, simulationConfig SimulationConfig, bettingStrategy bettingstrategy.Strategy, players []blackjack.BlackjackPlayer, bj *blackjack.Blackjack, i uint) error {
	logger.Debug("starting round", zap.Uint("round", i))

	if players[0].Credits == 0 {
		roundsSinceLastCredit := i - s.lastCreditRound

		if simulationConfig.BankCredits >= simulationConfig.StartingCredits {
			players[0].Credits = simulationConfig.StartingCredits
			simulationConfig.BankCredits -= simulationConfig.StartingCredits
			s.numberOfWithdrawals++
			logger.Debug("withdrew", zap.Uint64("credits", simulationConfig.StartingCredits), zap.Uint("round", i))
		} else {
			players[0].Credits = simulationConfig.BankCredits
			simulationConfig.BankCredits = 0
			logger.Info("bank is out of credits", zap.Uint("round", i))
			s.numberOfWithdrawals++
			logger.Debug("withdrew", zap.Uint64("credits", players[0].Credits), zap.Uint("round", i))
		}

		players[0].Hands[0].BetAmount = simulationConfig.OneCreditAmount
		simulationConfig.RebuyCount--
		logger.Debug("player rebuy", zap.Int("rebuys remaining", simulationConfig.RebuyCount))

		s.lastCreditRound = i
		s.creditAtRound = append(s.creditAtRound, roundsSinceLastCredit)
	}

	if players[0].Credits >= simulationConfig.BankAtCredits {
		simulationConfig.BankCredits += players[0].Credits - simulationConfig.StartingCredits
		s.numberOfDeposits++
		logger.Debug("deposited", zap.Uint64("credits", players[0].Credits-simulationConfig.StartingCredits), zap.Uint("round", i))
		players[0].Credits = simulationConfig.StartingCredits
	}

	profitPercentage := float64(players[0].Credits) / float64(simulationConfig.StartingCredits)
	if profitPercentage > s.highestProfitPercentage {
		s.highestProfitPercentage = profitPercentage
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

	return nil
}

func initTestPlayers(simulationConfig SimulationConfig) []blackjack.BlackjackPlayer {
	return []blackjack.BlackjackPlayer{
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
}

func (s *Simulator) logSimulationCompletion() {
	s.logger.Info("simulation complete",
		zap.Uint("rounds", s.currentRound),
		zap.Int("remaining rebuys", s.startingRebuyCount-s.RebuyCount),
		zap.Uint("deposits", s.numberOfDeposits),
		zap.Uint("withdrawals", s.numberOfWithdrawals),
		zap.String("deposit percentage", fmt.Sprintf("%.2f%%", s.getDepositPercentage()*100)))
}

func (s *Simulator) getDepositPercentage() float64 {
	if s.numberOfWithdrawals == 0 {
		return 1
	}

	return float64(s.numberOfDeposits) / float64(s.numberOfWithdrawals)
}
