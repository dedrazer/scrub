package models

type SimulationResults struct {
	AverageRoundsSurvived      uint    `json:"average_rounds_survived"`
	EarliestBankruptcyRound    uint    `json:"earliest_bankruptcy_round"`
	HighestProfitPercentage    float64 `json:"highest_profit_percentage"`
	OneCreditPercentageOfTotal float64 `json:"one_credit_percentage_of_total"`
	Score                      int
}
