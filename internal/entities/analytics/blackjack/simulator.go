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
		bettingStrategy:         strategy,
		startingTime:            time.Now().UTC(),
		startingRebuyCount:      config.RebuyCount,
		startingBankCredits:     config.BankCredits,
		lastCreditRound:         uint(0),
		creditAtRound:           make([]uint, 1),
		highestProfitPercentage: float64(1),
		SimulationConfig:        config,
	}
}

func (s *Simulator) Simulate() error {
	s.logger.Info("starting simulation", zap.Any("config", s.SimulationConfig))

	s.players = initTestPlayers(s.SimulationConfig)

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
	s.logger.Info("runtime statistics",
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
	oneCreditPercentageOfTotal := float64(s.OneCreditAmount) / float64(s.StartingCredits)

	res := models.SimulationResults{
		AverageRoundsSurvived:      uint(averageRoundsSurvived),
		EarliestBankruptcyRound:    earliestBankruptcyRound,
		HighestProfitPercentage:    s.highestProfitPercentage,
		OneCreditPercentageOfTotal: oneCreditPercentageOfTotal,
		StartingCredits:            s.startingBankCredits,
		EndingCredits:              s.BankCredits,
		RebuyCredits:               s.StartingCredits,
		BankAtCredits:              s.BankAtCredits,
		Score:                      float64(s.highestProfitPercentage) * s.getDepositPercentage() * float64(averageRoundsSurvived) * float64(oneCreditPercentageOfTotal),
	}

	s.logger.Info("strategy results", zap.Any("results", res))

	return nil
}

func (s *Simulator) simulateRounds() error {
	for s.currentRound < s.MaxRounds && s.hasPositiveBalance() {
		err := s.simulateRound()
		if err != nil {
			return errors.ErrFailedSubMethod("simulateRound", err)
		}

		s.currentRound++
	}

	return nil
}

func (s *Simulator) simulateRound() error {
	s.logger.Debug("starting round", zap.Uint("round", s.currentRound))

	if s.players[0].Credits == 0 {
		roundsSinceLastCredit := s.currentRound - s.lastCreditRound

		s.withdrawFromBank()

		s.players[0].Hands[0].BetAmount = s.OneCreditAmount
		s.RebuyCount--
		s.logger.Debug("player rebuy", zap.Int("rebuys remaining", s.RebuyCount))

		s.lastCreditRound = s.currentRound
		s.creditAtRound = append(s.creditAtRound, roundsSinceLastCredit)
	}

	if s.players[0].Credits >= s.BankAtCredits {
		s.BankCredits += s.players[0].Credits - s.StartingCredits
		s.numberOfDeposits++
		s.logger.Debug("deposited", zap.Uint64("credits", s.players[0].Credits-s.StartingCredits), zap.Uint("round", s.currentRound))
		s.players[0].Credits = s.StartingCredits
	}

	profitPercentage := float64(s.players[0].Credits) / float64(s.StartingCredits)
	if profitPercentage > s.highestProfitPercentage {
		s.highestProfitPercentage = profitPercentage
	}

	err := s.bettingStrategy.Strategy(s.players)
	if err != nil {
		return errors.ErrFailedSubMethod("bettingStrategy", err)
	}

	if s.players[0].Hands[0].BetAmount > s.players[0].Credits {
		s.players[0].Hands[0].BetAmount = s.players[0].Credits
	}

	dealerHand, err := s.blackjackEngine.DealRound(s.players)
	if err != nil {
		return errors.ErrFailedSubMethod("DealRound", err)
	}

	err = s.blackjackEngine.Play(s.logger, s.players, dealerHand, blackjack.Strategy)
	if err != nil {
		return errors.ErrFailedSubMethod("Play", err)
	}

	return nil
}

func (s *Simulator) withdrawFromBank() {
	var amount uint64

	if s.BankCredits >= s.StartingCredits {
		s.players[0].Credits = s.StartingCredits
		s.BankCredits -= s.StartingCredits
		amount = s.StartingCredits
	} else {
		s.players[0].Credits = s.BankCredits
		s.BankCredits = 0
		s.logger.Info("bank is out of credits", zap.Uint("round", s.currentRound))
		amount = s.players[0].Credits
	}

	s.logger.Debug("withdrew", zap.Uint64("credits", amount), zap.Uint("round", s.currentRound))
	s.numberOfWithdrawals++
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

func (s *Simulator) logPlayersStatistics() {
	for j := range s.players {
		s.players[j].LogStatistics(s.logger)
	}
}

func (s *Simulator) hasPositiveBalance() bool {
	return s.players[0].Credits > 0 || (s.RebuyCount > 0 && s.BankCredits > 0)
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

	return fmt.Sprintf("%.2fh", float64(totalDurationMs)/3600000)
}

func (s *Simulator) getDepositPercentage() float64 {
	if s.numberOfWithdrawals == 0 {
		return 1
	}

	return float64(s.numberOfDeposits) / float64(s.numberOfWithdrawals)
}
