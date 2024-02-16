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
	id, err := service.ResolveId("https://tidal.com/track/634872 best trip hop album imo! with his other one maxinquaye")
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	metadata, err := linkRes.GetTrackMetadata(id)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	result, err := linkRes.GetLinksFromMetadata(metadata)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	log.Info().Any("result", result).Send()
}
