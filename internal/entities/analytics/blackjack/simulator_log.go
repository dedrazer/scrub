package blackjackanalytics

import (
	"fmt"

	"go.uber.org/zap"
)

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

func (s *Simulator) logRuntimeStatistics(totalDurationTextual, averageRoundDuration string, roundsPerSecond int64) {
	s.logger.Info("runtime statistics",
		zap.String("duration", totalDurationTextual),
		zap.Int64("rounds per second", roundsPerSecond),
		zap.String("average round duration", averageRoundDuration))
}
