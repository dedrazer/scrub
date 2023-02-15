package player

import (
	"fmt"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

type Player struct {
	Name            string
	StartingCredits uint64
	Credits         uint64
	Wins            uint64
	Losses          uint64
	Draws           uint64
}

func (p *Player) Win(amount uint64) {
	p.Credits += amount

	if amount > 0 {
		p.Wins++
	} else {
		p.Draws++
	}
}

func (p *Player) Lose(amount uint64) error {
	if p.Credits < amount {
		return errors.ErrInsufficientCredits
	}

	p.Credits -= amount

	p.Losses++
	return nil
}

func (p *Player) WinRate() float64 {
	return float64(p.Wins) / float64(p.Wins+p.Losses)
}

func (p *Player) WinRateString() string {
	return fmt.Sprintf("%.2f%%", p.WinRate()*100)
}

func (p *Player) LogStatistics(logger *zap.Logger) {
	logger.Info("player credit statistics",
		zap.Uint64("credits", p.Credits),
		zap.String("percentage", fmt.Sprintf("%.2f%%", (float64(p.Credits)/float64(p.StartingCredits))*100)))

	logger.Info("player match statistics",
		zap.Uint64("won", p.Wins),
		zap.Uint64("lost", p.Losses),
		zap.String("win rate", p.WinRateString()),
		zap.Uint64("drawn", p.Draws))
}
