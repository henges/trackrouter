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
	linkRes := service.NewLinkResolutionService(deps.Clients)
	result, err := linkRes.FindLinks("https://tidal.com/track/634872 best trip hop album imo! with his other one maxinquaye")
	if err != nil {
		log.Err(err).Send()
	}
	log.Info().Any("result", result).Send()
}
