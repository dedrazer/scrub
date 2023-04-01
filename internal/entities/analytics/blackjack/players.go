package blackjackanalytics

import (
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/player"
)

func getTestPlayers(simulationConfig SimulationConfig) []blackjack.BlackjackPlayer {
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
