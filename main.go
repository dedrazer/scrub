package main

import (
	"fmt"
	"scrub/internal/entities/blackjack"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %w ", err))
	}

	//deck.Demo(logger)
	blackjack.Demo(logger)
}
