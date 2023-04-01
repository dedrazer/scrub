package blackjackanalytics

import (
	"errors"
	"fmt"
	"scrub/internal/entities/analytics/blackjack/models"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/bettingstrategy"
	"scrub/internal/errorutils"
	"scrub/internal/utils"
	"time"

	"go.uber.org/zap"
)

type SimulationConfig struct {
	MaxRounds       uint
	Decks           uint
	StartingCredits uint64
	BankCredits     uint64
	BankAtCredits   uint64
	OneCreditAmount uint64
	RebuyCount      int
}

type Simulator struct {
	logger                  *zap.Logger
	bettingStrategy         bettingstrategy.Strategy
	blackjackEngine         *blackjack.Blackjack
	players                 []blackjack.BlackjackPlayer
	startingTime            time.Time
	startingRebuyCount      int
	startingBankCredits     uint64
	highestProfitPercentage float64
	numberOfDeposits        uint
	numberOfWithdrawals     uint
	currentRound            uint
	lastRebuyRound          uint
	creditAtRound           []uint
	earliestBankruptcyRound uint
	averageRoundsSurvived   float64
	SimulationConfig
}

func NewSimulator(logger *zap.Logger, strategy bettingstrategy.Strategy, config SimulationConfig) *Simulator {
	return &Simulator{
		logger:                  logger,
		bettingStrategy:         strategy,
		startingTime:            time.Now().UTC(),
		startingRebuyCount:      config.RebuyCount,
		startingBankCredits:     config.BankCredits,
		lastRebuyRound:          uint(0),
		creditAtRound:           make([]uint, 1),
		earliestBankruptcyRound: config.MaxRounds,
		highestProfitPercentage: float64(1),
		SimulationConfig:        config,
	}
}

func (s *Simulator) Simulate() error {
	s.logger.Info("starting simulation", zap.Any("config", s.SimulationConfig))

	s.players = getTestPlayers(s.SimulationConfig)

	s.blackjackEngine = blackjack.NewBlackjack(s.SimulationConfig.Decks)

	s.currentRound = 0

	err := s.simulateRounds()
	if err != nil {
		return err
	}

	s.logSimulationCompletion()

	s.logPlayersStatistics()

	s.blackjackEngine.LogStatistics(s.logger)

	totalDurationMs := time.Since(s.startingTime).Milliseconds()
	totalDurationTextual := s.getTextualDuration(totalDurationMs)

	averageRoundDuration := fmt.Sprintf("%.2fμs", (float64(totalDurationMs)/float64(s.MaxRounds))*1000)
	roundsPerSecond := int64(float64(s.MaxRounds*1000) / float64(totalDurationMs))

	s.logRuntimeStatistics(totalDurationTextual, averageRoundDuration, roundsPerSecond)

	s.averageRoundsSurvived = float64(s.currentRound) / float64(len(s.creditAtRound))

	res := s.getSimulationResults()

	s.logger.Info("strategy results", zap.Any("results", res))

	return nil
}

func (s *Simulator) simulateRounds() error {
	for s.currentRound < s.MaxRounds && s.hasRemainingBalance() {
		err := s.simulateRound()
		if err != nil {
			return errorutils.ErrFailedSubMethod("simulateRound", err)
		}

		s.currentRound++
	}

	return nil
}

func (s *Simulator) simulateRound() error {
	s.logger.Debug("starting round", zap.Uint("round", s.currentRound))

	if s.players[0].Credits == 0 {
		err := s.rebuy()
		if err != nil {
			return errorutils.ErrFailedSubMethod("rebuy", err)
		}
	}

	s.depositExcessIntoBank()

	s.highestProfitPercentage = s.recalculateHighestProfitPercentage()

	err := s.bettingStrategy.Strategy(s.players)
	if err != nil {
		return errorutils.ErrFailedSubMethod("bettingStrategy", err)
	}

	if s.players[0].Hands[0].BetAmount > s.players[0].Credits {
		s.players[0].Hands[0].BetAmount = s.players[0].Credits
	}

	dealerHand, err := s.blackjackEngine.DealRound(s.players)
	if err != nil {
		return errorutils.ErrFailedSubMethod("DealRound", err)
	}

	err = s.blackjackEngine.Play(s.logger, s.players, dealerHand, blackjack.Strategy)
	if err != nil {
		return errorutils.ErrFailedSubMethod("Play", err)
	}

	return nil
}

func (s *Simulator) rebuy() error {
	if s.RebuyCount < 1 {
		return errors.New("no rebuys remaining")
	}

	roundsSinceLastRebuy := s.currentRound - s.lastRebuyRound

	s.withdrawFromBank()

	s.players[0].Hands[0].BetAmount = s.OneCreditAmount
	s.RebuyCount--
	s.logger.Debug("player rebuy", zap.Int("rebuys remaining", s.RebuyCount))

	s.lastRebuyRound = s.currentRound
	s.creditAtRound = append(s.creditAtRound, roundsSinceLastRebuy)
	if roundsSinceLastRebuy < s.earliestBankruptcyRound {
		s.earliestBankruptcyRound = roundsSinceLastRebuy
	}

	return nil
}

func (s *Simulator) withdrawFromBank() {
	var amount uint64

	if s.BankCredits >= s.StartingCredits {
		amount = s.StartingCredits
	} else {
		amount = s.BankCredits
		s.logger.Info("bank is out of credits", zap.Uint("round", s.currentRound))
	}

	s.players[0].Credits += amount
	s.BankCredits -= amount

	s.logger.Debug("withdrew", zap.Uint64("credits", amount), zap.Uint("round", s.currentRound))
	s.numberOfWithdrawals++
}

func (s *Simulator) depositExcessIntoBank() {
	if s.players[0].Credits >= s.BankAtCredits {
		s.BankCredits += s.players[0].Credits - s.StartingCredits
		s.numberOfDeposits++
		s.logger.Debug("deposited", zap.Uint64("credits", s.players[0].Credits-s.StartingCredits), zap.Uint("round", s.currentRound))
		s.players[0].Credits = s.StartingCredits
	}
}

func (s *Simulator) recalculateHighestProfitPercentage() float64 {
	profitPercentage := float64(s.players[0].Credits) / float64(s.StartingCredits)
	if profitPercentage > s.highestProfitPercentage {
		return profitPercentage
	}

	return s.highestProfitPercentage
}

func (s *Simulator) getSimulationResults() models.SimulationResults {
	oneCreditPercentageOfTotal := s.getOneCreditPercentageOfStartingCredits()
	return models.SimulationResults{
		AverageRoundsSurvived:      uint(s.averageRoundsSurvived),
		EarliestBankruptcyRound:    s.earliestBankruptcyRound,
		HighestProfitPercentage:    s.highestProfitPercentage,
		OneCreditPercentageOfTotal: oneCreditPercentageOfTotal,
		StartingCredits:            s.startingBankCredits,
		EndingCredits:              s.BankCredits,
		RebuyCredits:               s.StartingCredits,
		BankAtCredits:              s.BankAtCredits,
		Score:                      s.getScore(),
	}
}

func (s *Simulator) getTextualDuration(totalDurationMs int64) string {
	if totalDurationMs < 1000 {
		return fmt.Sprintf("%dms", totalDurationMs)
	}

	if totalDurationMs < 60000 {
		return fmt.Sprintf("%.2fsec", float64(totalDurationMs)/1000)
	}

	if totalDurationMs < 3600000 {
		return fmt.Sprintf("%.2fmin", float64(totalDurationMs)/60000)
	}

	return fmt.Sprintf("%.2fhrs", float64(totalDurationMs)/3600000)
}

func (s *Simulator) getScore() float64 {
	score := float64(s.highestProfitPercentage) * s.getDepositPercentage() * float64(s.averageRoundsSurvived) * s.getOneCreditPercentageOfStartingCredits()
	return utils.Round(score, 2)
}

func (s *Simulator) getDepositPercentage() float64 {
	if s.numberOfWithdrawals == 0 {
		return 1
	}

	return float64(s.numberOfDeposits) / float64(s.numberOfWithdrawals)
}

func (s *Simulator) getOneCreditPercentageOfStartingCredits() float64 {
	return float64(s.OneCreditAmount) / float64(s.StartingCredits)
}

func (s *Simulator) hasRemainingBalance() bool {
	return s.players[0].Credits > 0 || (s.RebuyCount > 0 && s.BankCredits > 0)
}
