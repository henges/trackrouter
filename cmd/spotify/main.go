package main

import (
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/di"
	"github.com/rs/zerolog/log"
)

func main() {

	c := config.Get()
	deps := di.Get(c)

	token, err := deps.Clients.SpotifyClient.Token()
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	log.Info().Any("token", token).Send()
}
