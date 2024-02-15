package main

import (
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/di"
	"github.com/henges/trackrouter/service"
	"github.com/rs/zerolog/log"
)

func main() {

	c := config.Get()
	deps := di.Get(c)
	linkRes := service.NewLinkResolutionService(deps)

	result, err := linkRes.GetLinks("dean blunt narcissist")
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	log.Info().Any("result", result).Send()
}
